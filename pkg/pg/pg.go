package pg

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx"
	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	queryModifierRegex = `^select\s+([a-zA-Z,\s_]*)\s+from\s+([a-zA-Z,_]*)(\s*where\s+[a-zA-Z0-9_\s><!()\[\]'"]*)*(\s+order\s+[a-zA-Z_\s]*)*(;?$)`
	paginationQuery    = `with filtered as (select $1 from $2$3$4), selected as (select *, count(*) over () as %s, row_number() over () as %s from filtered) select * from selected where %s > $%d limit $%d;`
)

type Pg struct {
	QueryInterface        *pgx.ConnPool
	NotificationInterface *pgx.Conn
	config                pgx.ConnConfig
}

func (pg *Pg) Close() error {
	errs := make([]string, 0)
	pg.QueryInterface.Close()
	if pg.NotificationInterface != nil {
		err := pg.NotificationInterface.Close()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.Errorf(
			`Encountered following errors while closing connection: %s`,
			strings.Join(errs, `, `),
		)
	}

	return nil
}

func (pg *Pg) GetObject(query types.Query) (map[string]interface{}, error) {
	transformed := make(map[string]interface{}, 1)
	result, err := pg.QueryInterface.Query(query.QueryString, query.Args...)
	if err != nil {
		return nil, err
	}

	if result.Next() {
		descriptors := result.FieldDescriptions()
		values, err := result.Values()
		if err != nil {
			transformed = nil
		} else {
			for i := range descriptors {
				val := SafeParseValue(values[i])
				transformed[descriptors[i].Name] = val
			}
			err = nil
		}

	} else {
		transformed = nil
		err = errors.NotFoundf(`record not found`)
	}
	result.Close()
	return transformed, err
}

func (pg *Pg) GetArray(query types.Query) ([]map[string]interface{}, error) {
	transformed := make([]map[string]interface{}, 0)
	result, err := pg.QueryInterface.Query(query.QueryString, query.Args...)

	if err != nil {
		transformed = nil
	} else {
		for result.Next() {
			descriptors := result.FieldDescriptions()
			values, err := result.Values()
			if err != nil {
				continue
			}
			each := make(map[string]interface{}, 1)
			for i := range descriptors {
				each[descriptors[i].Name] = values[i]
			}
			transformed = append(transformed, each)
		}
		err = nil

	}
	return transformed, err
}

func (pg *Pg) GetPaginatedArray(query types.Query, previous types.Pagination, recordIdKey string) ([]map[string]interface{}, *types.Pagination, error) {
	queryModifier := regexp.MustCompile(queryModifierRegex)
	argSize := len(query.Args)
	partial := queryModifier.ReplaceAllString(
		query.QueryString,
		paginationQuery,
	)
	modifiedQuery := fmt.Sprintf(
		partial,
		types.QueryTotalCount,
		types.QueryPaginationId,
		types.QueryPaginationId,
		argSize+1,
		argSize+2,
	)

	transformed := make([]map[string]interface{}, 0)
	pagination := &types.Pagination{}
	result, err := pg.QueryInterface.Query(modifiedQuery, append(query.Args, previous.PaginationId, previous.Limit)...)

	if err != nil {
		transformed = nil
		pagination = nil
	} else {
		recordCount := 0
		for result.Next() {
			descriptors := result.FieldDescriptions()
			values, err := result.Values()
			if err != nil {
				continue
			}

			pagination.Limit += 1
			each := make(map[string]interface{}, 1)
			for i := range descriptors {
				key := descriptors[i].Name
				value := values[i]
				switch key {
				case types.QueryTotalCount:
					pagination.TotalCount = value.(int64)
					break
				case types.QueryPaginationId:
					pagination.PaginationId = value.(int64)
					break
				default:
					each[key] = values[i]
				}
				if key == recordIdKey {
					parsed := value.(int64)
					if recordCount == 0 {
						pagination.LowestRecordId = parsed
						pagination.HighestRecordId = parsed
					}
					if pagination.LowestRecordId > parsed {
						pagination.LowestRecordId = parsed
					}
					if pagination.HighestRecordId < parsed {
						pagination.HighestRecordId = parsed
					}
				}
			}
			recordCount += 1
			transformed = append(transformed, each)
			err = nil
		}
	}

	return transformed, pagination, err
}

func (pg *Pg) Set(queries []types.Query) error {
	tx, err := pg.QueryInterface.Begin()
	if err != nil {
		return err
	}

	//noinspection GoUnhandledErrorResult
	defer tx.Rollback()

	for _, query := range queries {
		_, err = tx.Exec(query.QueryString, query.Args...)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (pg *Pg) Notify(listener types.NotificationListener) error {
	if pg.NotificationInterface == nil {
		conn, err := pgx.Connect(pg.config)
		if err != nil {
			return err
		}
		pg.NotificationInterface = conn
	}
	go func() {
		for pg.NotificationInterface.IsAlive() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			notification, err := pg.NotificationInterface.WaitForNotification(ctx)
			cancel()

			if err == nil && notification != nil {
				listener.OnNotification(NewNotification(notification))
			} else {
				time.Sleep(10 * time.Second)
			}
		}
	}()
	return nil
}

func (pg *Pg) GenerateUpdateQuery(tableName string, args map[string]interface{}, constraints []types.Constraint) types.Query {
	argCount := 0

	updateColumns := make([]string, 0)
	updatedColumnValues := make([]interface{}, 0)

	for key, value := range args {
		argCount += 1
		updateColumns = append(
			updateColumns,
			fmt.Sprintf(`%s = $%d`, key, argCount),
		)

		updatedColumnValues = append(updatedColumnValues, value)
	}

	whereColumns := make([]string, 0)
	whereColumnValues := make([]interface{}, 0)

	for _, constraint := range constraints {
		argCount += 1
		whereColumns = append(
			whereColumns,
			fmt.Sprintf(`%s %s $%d`, constraint.Field, constraint.Relation, argCount),
		)
		whereColumnValues = append(whereColumnValues, constraint.Value)
	}

	return types.Query{
		QueryString: fmt.Sprintf(
			`update %s set %s where %s;`,
			tableName,
			strings.Join(updateColumns, `, `),
			strings.Join(whereColumns, `, `),
		),
		Args: append(updatedColumnValues, whereColumnValues...),
	}
}

func (pg *Pg) Listen(channels ...string) error {
	if pg.NotificationInterface == nil {
		return errors.NotProvisionedf(`connection for listening to notification not established`)
	}
	errorChannels := make([]string, 0)
	for _, channel := range channels {
		err := pg.NotificationInterface.Listen(channel)
		if err != nil {
			errorChannels = append(errorChannels, channel)
		}
	}
	if len(errorChannels) > 0 {
		return errors.Errorf(`UNABLE TO CONNECT TO  CHANNEL(S) %s`, strings.Join(errorChannels, `, `))
	}
	return nil
}
