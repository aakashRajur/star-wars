package module

import (
	"context"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/redis"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetRedis(lifecycle fx.Lifecycle, logger types.Logger, handler types.FatalHandler) *redis.Redis {
	redisUri := env.GetString("REDIS_URI")
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
