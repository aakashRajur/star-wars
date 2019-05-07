package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/consul"
	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/resource"
)

func GetConsul() *consul.Consul {
	consulHost := env.GetString(`SERVICE_DISCOVERY_URI`)
	instance := consul.Consul{Host: consulHost}

	return &instance
}

func GetResourceRegistration(consul *consul.Consul) resource.Registree {
	return consul
}

var ConsulModule = fx.Provide(
	GetConsul,
	GetResourceRegistration,
)
