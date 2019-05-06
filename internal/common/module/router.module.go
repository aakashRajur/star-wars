package module

import (
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
	"go.uber.org/fx"
)

func GetHttpRouter(logger types.Logger, resources []http.Resource) *http.Router {
	router := http.NewRouter(logger)
	for _, each := range resources {
		router.DefineResource(each)
	}

	return router
}

var HttpRouterModule = fx.Provide(GetHttpRouter)
