package vehicles

import (
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/redis-pg"
)

const (
	CacheKey = `vehicles`
)

func CacheStrategy() di.CacheStrategyProvider {
	strategy := redis_pg.CacheStrategy{
		Channel:  CacheKey,
		CacheKey: CacheKey,
	}

	return di.CacheStrategyProvider{
		Strategy: strategy,
	}
}
