package planet

import (
	nativeHttp "net/http"

	"github.com/aakashRajur/star-wars/internal/wiki/api/planet"
	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func PatchPlanet(storage types.Storage, logger types.Logger, tracker types.TimeTracker, paramKey string) http.HandlerWithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		params := request.GetParams()

		id := params[paramKey].(int)

		ctx := request.Context()
		body := ctx.Value(middleware.JSON_BODY).(map[string]interface{})

		err := planet.QueryUpdatePlanet(storage, tracker, id, body)

		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		response.WriteHeader(nativeHttp.StatusOK)
	}

	middlewares := http.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.ValidateArgs(
			planet.ArgValidation,
			planet.ArgNormalization,
		),
		middleware.JsonBodyParser,
		middleware.ValidateBody(
			planet.BodyValidation,
			planet.BodyNormalization,
		),
	)

	return http.HandlerWithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
