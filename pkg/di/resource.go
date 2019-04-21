package di

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/http"
)

type ResourcesCompiler struct {
	fx.In
	Resources []http.Resource `group:"resources"`
}

type ResourceProvider struct {
	fx.Out
	Resource http.Resource `group:"resources"`
}
