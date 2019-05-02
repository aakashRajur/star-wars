package kafka

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/module"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetKafkaHealthCheck(lifecycle fx.Lifecycle, logger types.Logger, resources []http.Resource, handler types.FatalHandler) {
	server := module.GetHttpServer(logger, resources, handler)
	lifecycle.Append(
		fx.Hook{
			OnStart: server.Start,
			OnStop:  server.Stop,
		},
	)
}

var HealthcheckModule = fx.Invoke(GetKafkaHealthCheck)
