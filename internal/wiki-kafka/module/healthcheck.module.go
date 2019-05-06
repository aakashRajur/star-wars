package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/http"
)

func GetKafkaHealthCheck(lifecycle fx.Lifecycle, server *http.Server) {
	lifecycle.Append(
		fx.Hook{
			OnStart: server.Start,
			OnStop:  server.Stop,
		},
	)
}

var HealthcheckModule = fx.Invoke(GetKafkaHealthCheck)
