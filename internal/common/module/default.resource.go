package module

import (
	"github.com/aakashRajur/star-wars/pkg/service"
	"go.uber.org/fx"
)

func GetServiceResources() []service.Resource {
	return []service.Resource{}
}

var ResourceModule = fx.Provide(GetServiceResources)
