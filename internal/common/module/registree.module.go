package module

import (
	"context"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/service"
)

func RegisterService(lifecycle fx.Lifecycle, service service.Service, registree service.Registree) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				return registree.Register(service)
			},
			OnStop: func(context.Context) error {
				return registree.Unregister(service)
			},
		},
	)
}

var RegistreeModule = fx.Invoke(RegisterService)
