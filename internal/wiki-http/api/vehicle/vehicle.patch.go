package vehicle

import (
	nativeHttp "net/http"

	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicle"
	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func ApiPatchVehicle(storage types.Storage, logger types.Logger, tracker types.TimeTracker, paramKey string) http.HandlerWithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		params := request.GetParams()

		id := params[paramKey].(int)

		ctx := request.Context()
		body := ctx.Value(middleware.JSON_BODY).(map[string]interface{})

		err := vehicle.QueryUpdateVehicle(storage, tracker, id, body)

		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		response.WriteHeader(nativeHttp.StatusOK)
	}

	middlewares := http.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.ValidateArgs(
			vehicle.ArgValidation,
			vehicle.ArgNormalization,
		),
		middleware.JsonBodyParser,
		middleware.ValidateBody(
			vehicle.BodyValidation,
			vehicle.BodyNormalization,
		),
	)

	return http.HandlerWithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
