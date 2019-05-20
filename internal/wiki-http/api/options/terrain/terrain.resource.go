package terrain

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/options/terrain"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(terrain.HttpUrl)
	resource.Get(GetTerrain(storage, logger, tracker))

	return http_resource.HttpResourceProvider{
		Resource: resource,
	}
}
