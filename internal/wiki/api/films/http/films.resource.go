package films

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/films"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) di.ResourceProvider {
	resource := http.NewResource(films.HttpURI)
	resource.Get(GetFilms(storage, logger, tracker, films.CacheKey))

	return di.ResourceProvider{
		Resource: resource,
	}
}
