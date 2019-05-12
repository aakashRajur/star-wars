package di

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/http"
)

type HttpResourcesCompiler struct {
	fx.In
	Resources []http.Resource `group:"resources"`
}

type HttpResourceProvider struct {
	fx.Out
	Resource http.Resource `group:"resources"`
}
