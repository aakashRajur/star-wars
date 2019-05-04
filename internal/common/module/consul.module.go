package module

import (
	"github.com/aakashRajur/star-wars/pkg/consul"
	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/types"
	"go.uber.org/fx"
)

func GetConsul() *consul.Consul {
	consulHost := env.GetString(`SERVICE_DISCOVERY_URI`)
	instance := consul.Consul{Host: consulHost}

	return &instance
}

func GetResourceRegistration(consul *consul.Consul) types.ResourceRegistration {
	return consul
}

var ConsulModule = fx.Provide(
	GetConsul,
	GetResourceRegistration,
)
