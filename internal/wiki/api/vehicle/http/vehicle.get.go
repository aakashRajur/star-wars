package vehicle

import (
	nativeHttp "net/http"
	"strconv"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicle"
	middleware "github.com/aakashRajur/star-wars/middleware/http"

	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func ApiGetVehicle(storage types.Storage, logger types.Logger, tracker types.TimeTracker, cacheKey string, paramKey string) http.WithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		params := request.GetParams()

		id, ok := params[paramKey]
		if !ok {
			response.Error(nativeHttp.StatusNotAcceptable, errors.New(`vehicle id not provided`))
			return
		}
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			response.Error(nativeHttp.StatusUnprocessableEntity, errors.New(`invalid vehicle id`))
			return
		}

		data, err := vehicle.QuerySelectVehicle(storage, tracker, cacheKey, parsedId)
		if err != nil {
			response.Error(nativeHttp.StatusNotFound, err)
			return
		}

		err = response.WriteJSONObject(data, nil)
		if err != nil {
			logger.Error(err)
		}
	}

	middlewares := http.ChainMiddlewares(middleware.Logger(logger))

	return http.WithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
