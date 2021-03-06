package specie

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/specie"
	"github.com/aakashRajur/star-wars/internal/wiki/api/species"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(specie.HttpURL)
	resource.Get(ApiGetSpecie(storage, logger, tracker, species.CacheKey, specie.ParamSpecieId))
	resource.Patch(PatchSpecie(storage, logger, tracker, specie.ParamSpecieId))

	return http_resource.HttpResourceProvider{Resource: resource}
}
