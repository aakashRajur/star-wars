package planets

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/planets"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) di.ResourceProvider {
	resource := http.NewResource(planets.HttpURI)
	resource.Get(GetPlanets(storage, logger, tracker, planets.CacheKey))

	return di.ResourceProvider{Resource: resource}
}
