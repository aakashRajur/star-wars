package service_resource

import (
	"github.com/aakashRajur/star-wars/pkg/service"
	"go.uber.org/fx"
)

type ServiceResourceCompiler struct {
	fx.In
	Resources []service.Resource `group:"service_resources"`
}

type ServiceResourceProvider struct {
	fx.Out
	Resource service.Resource `group:"service_resources"`
}
