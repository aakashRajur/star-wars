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

func ApiPatchVehicle(storage types.Storage, logger types.Logger, tracker types.TimeTracker, paramKey string) http.HandlerWithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		params := request.GetParams()

		id, ok := params[paramKey]
		if !ok {
			response.Error(nativeHttp.StatusNotAcceptable, errors.New(`Planet id not provided`))
			return
		}
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			response.Error(nativeHttp.StatusUnprocessableEntity, errors.New(`Invalid planet id`))
			return
		}

		ctx := request.Context()
		body := ctx.Value(middleware.JSON_BODY).(map[string]interface{})

		err = vehicle.QueryUpdateVehicle(storage, tracker, parsedId, body)

		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		response.WriteHeader(nativeHttp.StatusOK)
	}

	middlewares := http.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.JsonBodyParser,
		middleware.ValidateBody(vehicle.VehicleValidation, vehicle.VehicleNormalization),
	)

	return http.HandlerWithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
