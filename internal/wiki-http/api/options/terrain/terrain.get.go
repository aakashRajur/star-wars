package terrain

import (
	nativeHttp "net/http"

	"github.com/aakashRajur/star-wars/internal/wiki/api/options/terrain"
	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetTerrain(storage types.Storage, logger types.Logger, tracker types.TimeTracker) http.HandlerWithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		data, err := terrain.QueryTerrain(storage, tracker)
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		err = response.WriteJSON(data, nil)
		if err != nil {
			logger.Error(err)
		}
	}

	middlewares := http.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.Pagination,
	)

	return http.HandlerWithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
