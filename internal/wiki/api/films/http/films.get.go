package films

import (
	"encoding/json"
	nativeHttp "net/http"

	"github.com/aakashRajur/star-wars/internal/wiki/api/films"
	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetFilms(storage types.Storage, logger types.Logger, tracker types.TimeTracker, cacheKey string) http.WithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		ctx := request.Context()
		oldPagination := ctx.Value(middleware.PAGINATION).(types.Pagination)
		result, newPagination, err := films.QuerySelectFilms(
			storage,
			tracker,
			cacheKey,
			oldPagination,
		)
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}
		marshaled, err := json.Marshal(*newPagination)
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		headers := make(map[string]string, 1)
		headers[middleware.PAGINATION] = string(marshaled)

		err = response.WriteJSONObject(result, &headers)
		if err != nil {
			logger.Error(err)
		}
	}
	middlewares := http.ChainMiddlewares(middleware.Logger(logger), middleware.Pagination)

	return http.WithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
