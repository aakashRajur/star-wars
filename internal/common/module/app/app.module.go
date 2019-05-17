package app

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func App(lifecycle fx.Lifecycle, protocol types.Protocol) {
	lifecycle.Append(
		fx.Hook{
			OnStart: protocol.Start,
			OnStop:  protocol.Stop,
		},
	)
}

var Module = fx.Invoke(App)
