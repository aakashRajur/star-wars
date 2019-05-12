package stats

import (
	http2 "net/http"

	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/util"
)

func Resource(logger types.Logger) di.HttpResourceProvider {
	resource := http.NewResource(`/stats`)

	resource.Get(
		http.HandlerWithMiddleware{
			Middlewares: func(next http.HandleRequest) http.HandleRequest {
				return func(response http.Response, request *http.Request) {
					util.LogRequest(logger, request)
					next(response, request)
				}
			},
			HandleRequest: func(response http.Response, request *http.Request) {
				stats, err := QueryStats()
				if err != nil {
					logger.Error(err)
					response.Error(http2.StatusInternalServerError, err)
					return
				}
				err = response.WriteJSON(stats, nil)
				if err != nil {
					logger.Error(err)
				}
			},
		},
	)

	return di.HttpResourceProvider{Resource: resource}
}
