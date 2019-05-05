package planet

import (
	nativeHttp "net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/internal/wiki/api/planet"
	middleware "github.com/aakashRajur/star-wars/middleware/http"

	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetPlanet(storage types.Storage, logger types.Logger, tracker types.TimeTracker, cacheKey string, paramKey string) http.HandlerWithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		params := request.GetParams()

		id, ok := params[paramKey]
		if !ok {
			response.Error(nativeHttp.StatusNotAcceptable, errors.New(`planet id not provided`))
			return
		}
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			response.Error(nativeHttp.StatusUnprocessableEntity, errors.New(`invalid planet id`))
			return
		}

		data, err := planet.QuerySelectPlanet(storage, tracker, cacheKey, parsedId)
		if err != nil {
			response.Error(nativeHttp.StatusNotFound, err)
			return
		}

		err = response.WriteJSON(data, nil)
		if err != nil {
			logger.Error(err)
		}
	}

	middlewares := http.ChainMiddlewares(middleware.Logger(logger))

	return http.HandlerWithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
