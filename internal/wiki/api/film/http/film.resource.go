package film

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/film"
	"github.com/aakashRajur/star-wars/internal/wiki/api/films"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) di.ResourceProvider {
	resource := http.NewResource(film.HttpURI)
	resource.Get(GetFilm(storage, logger, tracker, films.CacheKey, film.ParamFilmId))
	resource.Patch(PatchFilm(storage, logger, tracker, film.ParamFilmId))

	return di.ResourceProvider{Resource: resource}
}
