package pg

import (
	"github.com/jackc/pgx"
	"github.com/juju/errors"
)

func NewInstance(config Config, logger pgx.Logger) (*Pg, error) {
	pgUri := string(config.URL)
	if pgUri == "" {
		return nil, errors.NotValidf("PG_URI is invalid")
	}
	pgConfig, err := pgx.ParseConnectionString(pgUri)
	if err != nil {
		return nil, errors.NotSupportedf("DB_URL parse err")
	}
	pgConfig.Logger = logger

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     pgConfig,
		MaxConnections: config.PoolLimit,
	}

	queryInterface, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		return nil, errors.Annotate(err, "Postgres connection failed")
	}

	return &Pg{
		config:         pgConfig,
		QueryInterface: queryInterface,
	}, nil
}
