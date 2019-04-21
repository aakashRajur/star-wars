package http

import (
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Logger(logger types.Logger) http.Middleware {
	return func(next http.HandleRequest) http.HandleRequest {
		return func(response http.Response, request *http.Request) {
			params := request.Context().Value(http.PARAMS).(map[string]string)
			query := request.URL.Query()

			logger.InfoFields(
				map[string]interface{}{
					"METHOD":     request.Method,
					"URI":        request.RequestURI,
					"URL":        request.URL.Path,
					"PARAMS":     params,
					"QUERY":      query,
					"REMOTE":     request.RemoteAddr,
					"REFERER":    request.Referer(),
					"USER-AGENT": request.UserAgent(),
				},
			)

			next(response, request)
		}
	}
}
