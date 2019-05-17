package vehicle

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicle"
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicles"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(vehicle.HttpURL)
	resource.Get(ApiGetVehicle(storage, logger, tracker, vehicles.CacheKey, vehicle.ParamVehicleId))
	resource.Patch(ApiPatchVehicle(storage, logger, tracker, vehicle.ParamVehicleId))

	return http_resource.HttpResourceProvider{
		Resource: resource,
	}
}
