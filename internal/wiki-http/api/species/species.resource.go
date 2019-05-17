package species

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/species"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(species.HttpURL)
	resource.Get(GetSpecies(storage, logger, tracker, species.CacheKey))

	return http_resource.HttpResourceProvider{Resource: resource}
}
