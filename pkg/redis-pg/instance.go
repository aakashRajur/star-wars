package redis_pg

import (
	"github.com/aakashRajur/star-wars/pkg/pg"
	"github.com/aakashRajur/star-wars/pkg/redis"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func NewInstance(redis *redis.Redis, pg *pg.Pg, cacheStrategy []CacheStrategy, logger types.Logger) (*RedisPg, error) {
	instance := &RedisPg{
		pg:              pg,
		redis:           redis,
		updateListeners: make(map[string][]types.UpdateListener, 1),
		logger:          logger,
	}

	err := pg.Notify(instance)
	if err != nil {
		return nil, err
	}

	channels := make([]string, 0)
	for _, strategy := range cacheStrategy {
		channels = append(channels, strategy.Channel)
	}

	err = instance.Listen(CacheBuster(cacheStrategy), channels...)

	return instance, nil
}
