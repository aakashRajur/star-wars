package vehicles

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicles"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) di.ResourceProvider {
	resource := http.NewResource(vehicles.HttpURI)
	resource.Get(GetVehicles(storage, logger, tracker, vehicles.CacheKey))

	return di.ResourceProvider{
		Resource: resource,
	}
}