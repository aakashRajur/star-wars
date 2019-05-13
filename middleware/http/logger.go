package http

import (
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Logger(logger types.Logger) http.Middleware {
	return func(next http.HandleRequest) http.HandleRequest {
		return func(response http.Response, request *http.Request) {
			params := request.Context().Value(http.PARAMS).(map[string]interface{})
			query := request.URL.Query()

			logger.InfoFields(
				map[string]interface{}{
					"method":     request.Method,
					"uri":        request.RequestURI,
					"url":        request.URL.Path,
					"params":     params,
					"query":      query,
					"remote":     request.RemoteAddr,
					"referer":    request.Referer(),
					"user-agent": request.UserAgent(),
				},
			)
			next(response, request)
		}
	}
}
