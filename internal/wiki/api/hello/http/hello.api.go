package hello

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/hello"
	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(logger types.Logger) di.ResourceProvider {
	requestHandler := func(response http.Response, request *http.Request) {
		_, err := response.WriteText(hello.SayHello())
		if err != nil {
			logger.Error(err)
		}
	}
	middlewares := http.ChainMiddlewares(middleware.Logger(logger))
	resource := http.NewResource(`^/hello$`)
	resource.Get(
		http.HandlerWithMiddleware{
			HandleRequest: requestHandler,
			Middlewares:   middlewares,
		},
	)
	return di.ResourceProvider{
		Resource: resource,
	}
}
