package characters

import (
	"encoding/json"
	nativeHttp "net/http"

	"github.com/aakashRajur/star-wars/internal/wiki/api/characters"
	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetCharacters(storage types.Storage, logger types.Logger, tracker types.TimeTracker, cacheKey string) http.HandlerWithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		ctx := request.Context()
		oldPagination := ctx.Value(types.PAGINATION).(types.Pagination)
		result, newPagination, err := characters.QuerySelectCharacters(storage, tracker, cacheKey, oldPagination)
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
		headers[types.PAGINATION] = string(marshaled)

		err = response.WriteJSON(result, headers)
		if err != nil {
			logger.Error(err)
		}
	}

	middlewares := http.ChainMiddlewares(middleware.Logger(logger), middleware.Pagination)

	return http.HandlerWithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
