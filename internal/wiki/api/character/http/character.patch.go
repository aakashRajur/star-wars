package character

import (
	nativeHttp "net/http"

	"github.com/aakashRajur/star-wars/internal/wiki/api/character"
	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func PatchCharacter(storage types.Storage, logger types.Logger, tracker types.TimeTracker, paramKey string) http.HandlerWithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		params := request.GetParams()

		id := params[paramKey].(int)

		ctx := request.Context()
		body := ctx.Value(middleware.JSON_BODY).(map[string]interface{})

		err := character.QueryUpdateCharacter(storage, tracker, id, body)

		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		response.WriteHeader(nativeHttp.StatusOK)
	}

	middlewares := http.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.ValidateArgs(
			character.ArgValidation,
			character.ArgNormalization,
		),
		middleware.JsonBodyParser,
		middleware.ValidateBody(
			character.BodyValidation,
			character.BodyNormalization,
		),
	)

	return http.HandlerWithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
