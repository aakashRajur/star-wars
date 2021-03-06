package characters

import (
	"github.com/aakashRajur/star-wars/pkg/di/cache-strategy"
	"github.com/aakashRajur/star-wars/pkg/redis-pg"
)

const (
	CacheKey = `characters`
)

func CacheStrategy() cache_strategy.CacheStrategyProvider {
	strategy := redis_pg.CacheStrategy{
		Channel:  CacheKey,
		CacheKey: CacheKey,
	}

	return cache_strategy.CacheStrategyProvider{
		Strategy: strategy,
	}
}
