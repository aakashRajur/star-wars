package vehicle

import (
	nativeHttp "net/http"

	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicle"
	middleware "github.com/aakashRajur/star-wars/middleware/http"

	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func ApiGetVehicle(storage types.Storage, logger types.Logger, tracker types.TimeTracker, cacheKey string, paramKey string) http.HandlerWithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		params := request.GetParams()

		id := params[paramKey].(int)

		data, err := vehicle.QuerySelectVehicle(storage, tracker, cacheKey, id)
		if err != nil {
			response.Error(nativeHttp.StatusNotFound, err)
			return
		}

		err = response.WriteJSON(data, nil)
		if err != nil {
			logger.Error(err)
		}
	}

	middlewares := http.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.ValidateArgs(
			vehicle.ArgValidation,
			vehicle.ArgNormalization,
		),
	)

	return http.HandlerWithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
