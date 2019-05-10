package module

import (
	"time"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/consul"
	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

//noinspection GoSnakeCaseUsage
const (
	CONSUL_SCHEME  = `CONSUL_SCHEME`
	CONSUL_ADDRESS = `CONSUL_ADDRESS`
)

func GetConsulConfig() consul.Config {
	consulScheme := env.GetString(CONSUL_SCHEME)
	consulAddress := env.GetString(CONSUL_ADDRESS)

	config := consul.Config{
		Scheme:       consulScheme,
		Address:      consulAddress,
		WatchTimeout: 5 * time.Second,
	}
	return config
}

func GetConsul(config consul.Config, logger types.Logger) *consul.Consul {
	return consul.NewInstance(config, logger)
}

func GetRegistree(consul *consul.Consul) service.Registree {
	return consul
}

func GetResolver(consul *consul.Consul) service.Resolver {
	return consul
}

var ConsulModule = fx.Provide(
	GetConsulConfig,
	GetConsul,
	GetRegistree,
	GetResolver,
)
