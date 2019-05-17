package planet

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/planet"
	"github.com/aakashRajur/star-wars/internal/wiki/api/planets"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(planet.HttpURL)
	resource.Get(GetPlanet(storage, logger, tracker, planets.CacheKey, planet.ParamPlanetId))
	resource.Patch(PatchPlanet(storage, logger, tracker, planet.ParamPlanetId))

	return http_resource.HttpResourceProvider{Resource: resource}
}
