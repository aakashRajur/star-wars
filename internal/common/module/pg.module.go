package module

import (
	"context"

	"github.com/jackc/pgx"
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/pg"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetPsql(lifecycle fx.Lifecycle, logger pgx.Logger, handler types.FatalHandler) *pg.Pg {
	pgUri := env.GetString("DATABASE_URI")
	pgPoolLimit := env.GetInt("PG_POOL_LIMIT")
	psql, err := pg.NewInstance(
		pgUri,
		pgPoolLimit,
		logger,
	)
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

var PgModule = fx.Provide(
	GetPsql,
	GetPgLogger,
)
