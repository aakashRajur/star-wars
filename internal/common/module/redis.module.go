package module

import (
	"context"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/juju/errors"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/redis"
	"github.com/aakashRajur/star-wars/pkg/types"
)

//noinspection GoSnakeCaseUsage
const (
	REDIS_SERVICE = `CACHE_SERVICE`
)

func GetRedis(resolver service.Resolver, handler types.FatalHandler, logger types.Logger, lifecycle fx.Lifecycle) *redis.Redis {
	redisService := env.GetString(REDIS_SERVICE)
	endpoints, err := resolver.Resolve(redisService)
	if err != nil {
		handler.HandleFatal(err)
		return nil
	}
	if len(endpoints) < 1 {
		handler.HandleFatal(errors.New(`NO PG SERVICE FOUND`))
		return nil
	}

	redisUri := redis.Url(endpoints[0])
	client, err := redis.NewInstance(redisUri, logger)

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

var RedisModule = fx.Provide(GetRedis)
