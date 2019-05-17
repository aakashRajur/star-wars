package stats

import (
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	httpNative "net/http"

	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(logger types.Logger) http_resource.HttpResourceProvider {
	resource := http.NewResource(`/stats`)

	resource.Get(
		http.HandlerWithMiddleware{
			Middlewares: middleware.Logger(logger),
			HandleRequest: func(response http.Response, request *http.Request) {
				stats, err := QueryStats()
				if err != nil {
					logger.Error(err)
					response.Error(httpNative.StatusInternalServerError, err)
					return
				}
				err = response.WriteJSON(stats, nil)
				if err != nil {
					logger.Error(err)
				}
			},
		},
	)

	return http_resource.HttpResourceProvider{Resource: resource}
}
