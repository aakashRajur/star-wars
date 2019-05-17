package redis_pg

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/pg"
	"github.com/aakashRajur/star-wars/pkg/redis"
	"github.com/aakashRajur/star-wars/pkg/redis-pg"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetRedisPg(redis *redis.Redis, psql *pg.Pg, cacheStrategies []redis_pg.CacheStrategy, logger types.Logger, handler types.FatalHandler) *redis_pg.RedisPg {
	instance, err := redis_pg.NewInstance(redis, psql, cacheStrategies, logger)
	if err != nil {
		handler.HandleFatal(err)
	}
	return instance
}

func GetStorage(redisPg *redis_pg.RedisPg) types.Storage {
	return redisPg
}

var Module = fx.Provide(
	GetRedisPg,
	GetStorage,
)
