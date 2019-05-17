package vehicles

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicles"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(vehicles.HttpURL)
	resource.Get(GetVehicles(storage, logger, tracker, vehicles.CacheKey))

	return http_resource.HttpResourceProvider{
		Resource: resource,
	}
}
