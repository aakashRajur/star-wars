package http

import (
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/util"
)

func Logger(logger types.Logger) http.Middleware {
	return func(next http.HandleRequest) http.HandleRequest {
		return func(response http.Response, request *http.Request) {
			util.LogRequest(logger, request)
			next(response, request)
		}
	}
}
