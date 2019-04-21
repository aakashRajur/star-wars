package module

import (
	"time"

	"github.com/juju/errors"
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetHttpServer(logger types.Logger, resources []http.Resource, handler types.FatalHandler) *http.Server {
	sslCert := env.GetString(`HTTP_SSL_CERT`)
	if sslCert == `` {
		handler.HandleFatal(errors.New(`HTTP_SSL_CERT not set in env`))
	}
	sslKey := env.GetString(`HTTP_SSL_KEY`)
	if sslKey == `` {
		handler.HandleFatal(errors.New(`HTTP_SSL_KEY not set in env`))
	}
	httpPort := ":" + env.GetString(`HTTP_PORT`)
	timeoutConfig := http.Timeout{
		Read:  time.Second * time.Duration(env.GetInt(`HTTP_READ_TIMEOUT`)),
		Write: time.Second * time.Duration(env.GetInt(`HTTP_WRITE_TIMEOUT`)),
		Idle:  time.Second * time.Duration(env.GetInt(`HTTP_IDLE_TIMEOUT`)),
	}

	router := http.NewRouter(logger)
	for _, each := range resources {
		router.DefineResource(each)
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

var HttpModule = fx.Provide(
	GetHttpServer,
	GetHttpProtocol,
)
