package hello

import (
	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(logger types.Logger) http_resource.HttpResourceProvider {
	requestHandler := func(response http.Response, request *http.Request) {
		_, err := response.WriteText(SayHello())
		if err != nil {
			logger.Error(err)
		}
	}
	middlewares := http.ChainMiddlewares(middleware.Logger(logger))
	resource := http.NewResource(`/hello`)
	resource.Get(
		http.HandlerWithMiddleware{
			HandleRequest: requestHandler,
			Middlewares:   middlewares,
		},
	)
	return http_resource.HttpResourceProvider{
		Resource: resource,
	}
}
