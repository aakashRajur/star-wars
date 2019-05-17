package pg

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/juju/errors"
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/pg"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

//noinspection GoSnakeCaseUsage
const (
	PG_SERVICE = `DATABASE_SERVICE`
)

func GetPgUrl(resolver service.Resolver, handler types.FatalHandler) pg.Url {
	pgService := env.GetString(PG_SERVICE)
	endpoints, err := resolver.Resolve(pgService)
	if err != nil {
		handler.HandleFatal(err)
		return ``
	}
	if len(endpoints) < 1 {
		handler.HandleFatal(errors.New(`NO PG SERVICE FOUND`))
		return ``
	}
	return pg.Url(endpoints[0])
}

func GetPgConfig(url pg.Url) pg.Config {
	pgPoolLimit := env.GetInt("PG_POOL_LIMIT")
	return pg.Config{
		URL:       url,
		PoolLimit: pgPoolLimit,
	}
}

func GetPg(config pg.Config, handler types.FatalHandler, logger pgx.Logger, lifecycle fx.Lifecycle) *pg.Pg {
	psql, err := pg.NewInstance(config, logger)
	if err != nil {
		handler.HandleFatal(err)
		return nil
	}

	lifecycle.Append(
		fx.Hook{
			OnStop: func(context.Context) error {
				return psql.Close()
			},
		},
	)

	return psql
}

func GetPgLogger(logger types.Logger) pgx.Logger {
	return pg.NewPgLogger(logger)
}

var Module = fx.Provide(
	GetPgLogger,
	GetPgUrl,
	GetPgConfig,
	GetPg,
)
