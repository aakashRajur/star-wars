package planet

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/planet"
	"github.com/aakashRajur/star-wars/internal/wiki/api/planets"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) di.HttpResourceProvider {
	resource := http.NewResource(planet.HttpURL)
	resource.Get(GetPlanet(storage, logger, tracker, planets.CacheKey, planet.ParamPlanetId))
	resource.Patch(PatchPlanet(storage, logger, tracker, planet.ParamPlanetId))

	return di.HttpResourceProvider{Resource: resource}
}
