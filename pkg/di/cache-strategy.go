package di

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/redis-pg"
)

type CacheStrategyCompiler struct {
	fx.In
	Strategies []redis_pg.CacheStrategy `group:"cache-strategies"`
}

type CacheStrategyProvider struct {
	fx.Out
	Strategy redis_pg.CacheStrategy `group:"cache-strategies"`
}
