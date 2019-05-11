package module

import (
	"context"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/service"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/redis"
	"github.com/aakashRajur/star-wars/pkg/types"
)

//noinspection GoSnakeCaseUsage
const (
	REDIS_SERVICE = `CACHE_SERVICE`
)

func GetRedisUrl(resolver service.Resolver, handler types.FatalHandler) redis.Url {
	redisService := env.GetString(REDIS_SERVICE)
	endpoints, err := resolver.Resolve(redisService)
	if err != nil {
		handler.HandleFatal(err)
		return ``
	}
	if len(endpoints) < 1 {
		handler.HandleFatal(errors.New(`NO REDIS SERVICE FOUND`))
		return ``
	}

	return redis.Url(endpoints[0])
}

func GetRedis(url redis.Url, handler types.FatalHandler, logger types.Logger, lifecycle fx.Lifecycle) *redis.Redis {
	client, err := redis.NewInstance(url, logger)
	if err != nil {
		handler.HandleFatal(err)
	}

	lifecycle.Append(
		fx.Hook{
			OnStop: func(context.Context) error {
				return client.Close()
			},
		},
	)

	return client
}

var RedisModule = fx.Provide(
	GetRedisUrl,
	GetRedis,
)
