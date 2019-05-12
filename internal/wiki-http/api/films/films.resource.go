package films

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/films"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) di.HttpResourceProvider {
	resource := http.NewResource(films.HttpURL)
	resource.Get(GetFilms(storage, logger, tracker, films.CacheKey))

	return di.HttpResourceProvider{
		Resource: resource,
	}
}
