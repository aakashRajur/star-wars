package module

import (
	"time"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetHttpServer(logger types.Logger, router *http.Router) *http.Server {
	sslCert := env.GetString(`HTTP_SSL_CERT`)
	sslKey := env.GetString(`HTTP_SSL_KEY`)
	httpPort := ":" + env.GetString(`HTTP_PORT`)
	timeoutConfig := http.Timeout{
		Read:  time.Second * time.Duration(env.GetInt(`HTTP_READ_TIMEOUT`)),
		Write: time.Second * time.Duration(env.GetInt(`HTTP_WRITE_TIMEOUT`)),
		Idle:  time.Second * time.Duration(env.GetInt(`HTTP_IDLE_TIMEOUT`)),
	}

	serverConfig := http.ServerConfig{
		Port:    httpPort,
		Timeout: timeoutConfig,
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

var HttpModule = fx.Options(
	HttpRouterModule,
	fx.Provide(GetHttpServer),
)

var HttpProtocolModule = fx.Provide(GetHttpProtocol)
