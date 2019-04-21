package pg

import (
	"github.com/jackc/pgx"
	"github.com/juju/errors"
)

func NewInstance(uri string, poolLimit int, logger pgx.Logger) (*Pg, error) {
	if uri == "" {
		return nil, errors.NotValidf("PG_URI is invalid")
	}
	config, err := pgx.ParseConnectionString(uri)
	if err != nil {
		return nil, errors.NotSupportedf("DB_URL parse err")
	}
	config.Logger = logger

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: poolLimit,
	}

	queryInterface, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		return nil, errors.Annotate(err, "Postgres connection failed")
	}

	return &Pg{
		config:         config,
		QueryInterface: queryInterface,
	}, nil
}
