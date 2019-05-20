package climate

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/options/climate"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(climate.HttpUrl)
	resource.Get(GetClimate(storage, logger, tracker))

	return http_resource.HttpResourceProvider{
		Resource: resource,
	}
}
