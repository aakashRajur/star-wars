package planets

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/planets"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(planets.HttpURL)
	resource.Get(GetPlanets(storage, logger, tracker, planets.CacheKey))

	return http_resource.HttpResourceProvider{Resource: resource}
}
