package http

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/module/router"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetHttpServer(logger types.Logger, router *http.Router) *http.Server {
	sslCert := env.GetString(`HTTP_SSL_CERT`)
	sslKey := env.GetString(`HTTP_SSL_KEY`)
	httpPort := ":" + env.GetString(`HTTP_PORT`)

	serverConfig := http.ServerConfig{
		Port:    httpPort,
		SslCert: sslCert,
		SslKey:  sslKey,
		Logger:  logger,
	}
	server := http.NewInstance(serverConfig, router)
	router.Health = server.GetStatus

	return server
}

func GetHttpProtocol(server *http.Server) types.Protocol {
	return server
}

var Module = fx.Options(
	router.Module,
	fx.Provide(GetHttpServer),
)

var ProtocolModule = fx.Provide(GetHttpProtocol)
