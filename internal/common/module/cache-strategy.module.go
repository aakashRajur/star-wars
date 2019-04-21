package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/wiki/api/characters"
	"github.com/aakashRajur/star-wars/internal/wiki/api/films"
	"github.com/aakashRajur/star-wars/internal/wiki/api/planets"
	"github.com/aakashRajur/star-wars/internal/wiki/api/species"
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicles"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/redis-pg"
)

func GetCacheStrategies(cacheStrategies di.CacheStrategyCompiler) []redis_pg.CacheStrategy {
	return cacheStrategies.Strategies
}

var CacheStrategyModule = fx.Provide(
	planets.CacheStrategy,
	species.CacheStrategy,
	vehicles.CacheStrategy,
	characters.CacheStrategy,
	films.CacheStrategy,
	GetCacheStrategies,
)
