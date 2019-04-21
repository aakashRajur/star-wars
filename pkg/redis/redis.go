package redis

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/go-redis/redis"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	pageKey                = `pagination:%s:%d:%d`
	pageKeyPartial         = `pagination:%s:%s`
	pageLowKey             = `pagination:%s:lowest_record_id`
	pageHighKey            = `pagination:%s:highest_record_id`
	pageIndex              = `%06d:%d:%d`
	pageIndexPartial       = `[%06d`
	pageSetDataKey         = `data`
	pageSetPageKey         = `pagination`
	pageSetLowestRecordId  = `lowest_record_id`
	pageSetHighestRecordId = `highest_record_id`
)

type Redis struct {
	*redis.Client
}

func (client *Redis) GetObject(key string) (map[string]interface{}, error) {
	value, err := client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(value), &parsed)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func (client *Redis) GetArray(key string) ([]map[string]interface{}, error) {
	value, err := client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	var parsed []map[string]interface{}
	err = json.Unmarshal([]byte(value), &parsed)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func (client *Redis) GetPaginatedArray(key string, pagination types.Pagination) ([]map[string]interface{}, *types.Pagination, error) {
	dataKey := fmt.Sprintf(
		pageKey,
		key,
		pagination.PaginationId,
		pagination.Limit,
	)

	dataValue, err := client.HGet(
		dataKey,
		pageSetDataKey,
	).Result()

	if err != nil {
		return nil, nil, err
	}

	pageValue, err := client.HGet(
		dataKey,
		pageSetPageKey,
	).Result()

	if err != nil {
		return nil, nil, err
	}

	var parsedData []map[string]interface{}
	err = json.Unmarshal([]byte(dataValue), &parsedData)
	if err != nil {
		return nil, nil, err
	}

	parsedPage, err := types.PaginationFromString(pageValue)

	if err != nil {
		return nil, nil, err
	}

	return parsedData, &parsedPage, nil
}

func (client *Redis) SetObject(key string, data map[string]interface{}) error {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = client.Client.Set(key, marshaled, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *Redis) SetArray(key string, data []map[string]interface{}) error {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = client.Client.Set(key, marshaled, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *Redis) SetPaginatedArray(key string, old types.Pagination, data []map[string]interface{}, new types.Pagination) error {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	marshaledPage, err := json.Marshal(new)
	if err != nil {
		return err
	}

	dataKey := fmt.Sprintf(
		pageKey,
		key,
		old.PaginationId,
		old.Limit,
	)

	lowKeyValue := fmt.Sprintf(
		pageIndex,
		new.LowestRecordId,
		old.PaginationId,
		old.Limit,
	)

	highKeyValue := fmt.Sprintf(
		pageIndex,
		new.HighestRecordId,
		old.PaginationId,
		old.Limit,
	)

	pipe := client.TxPipeline()

	pipe.HSet(
		dataKey,
		pageSetDataKey,
		marshaledData,
	)

	pipe.HSet(
		dataKey,
		pageSetPageKey,
		marshaledPage,
	)

	pipe.HSet(
		dataKey,
		pageSetLowestRecordId,
		lowKeyValue,
	)

	pipe.HSet(
		dataKey,
		pageSetHighestRecordId,
		highKeyValue,
	)

	pipe.ZAdd(
		fmt.Sprintf(pageLowKey, key),
		redis.Z{
			Score:  0,
			Member: lowKeyValue,
		},
	)

	pipe.ZAdd(
		fmt.Sprintf(pageHighKey, key),
		redis.Z{
			Score:  0,
			Member: highKeyValue,
		},
	)

	_, err = pipe.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (client *Redis) Delete(keys ...string) error {
	_, err := client.Unlink(keys...).Result()
	return err
}

func (client *Redis) DeletePagination(keys ...string) error {
	pipe := client.TxPipeline()

	for _, key := range keys {
		entityParser := regexp.MustCompile(`^pagination:([a-zA-Z0-9_]*):\d+:\d+$`)
		entity := entityParser.ReplaceAllString(key, `$1`)

		lowIndex, _ := client.HGet(
			key,
			pageSetLowestRecordId,
		).Result()

		highIndex, _ := client.HGet(
			key,
			pageSetHighestRecordId,
		).Result()

		lowIndexKey := fmt.Sprintf(pageLowKey, entity)
		highIndexKey := fmt.Sprintf(pageHighKey, entity)

		_, err := pipe.ZRem(
			lowIndexKey,
			lowIndex,
		).Result()
		if err != nil {
			return err
		}

		_, err = pipe.ZRem(
			highIndexKey,
			highIndex,
		).Result()
		if err != nil {
			return err
		}

		pipe.Unlink(key)
	}

	_, err := pipe.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (client *Redis) GeneratePaginationCacheKeysForId(key string, id int64) ([]string, error) {
	parser := regexp.MustCompile(`\d{0,6}:(\d*:\d*)`)
	compiled := make([]string, 0)

	lesserThan, err := client.ZRangeByLex(
		fmt.Sprintf(pageLowKey, key),
		redis.ZRangeBy{
			Min: `-`,
			Max: fmt.Sprintf(pageIndexPartial, id+1),
		},
	).Result()
	if err != nil {
		return nil, err
	}

	greaterThan, err := client.ZRangeByLex(
		fmt.Sprintf(pageHighKey, key),
		redis.ZRangeBy{
			Min: fmt.Sprintf(pageIndexPartial, id),
			Max: `+`,
		},
	).Result()
	if err != nil {
		return nil, err
	}

	hashTrack := make(map[string]int, 1)

	for _, each := range lesserThan {
		hashTrack[parser.ReplaceAllString(each, `$1`)] = 0
	}

	for _, each := range greaterThan {
		parsed := parser.ReplaceAllString(each, `$1`)
		_, ok := hashTrack[parsed]
		if ok {
			hashTrack[parsed] += 1
			compiled = append(
				compiled,
				fmt.Sprintf(
					pageKeyPartial,
					key,
					parsed,
				),
			)
		}
	}

	return compiled, nil
}

func IsNil(err error) bool {
	return redis.Nil == err
}
