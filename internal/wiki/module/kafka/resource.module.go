package kafka

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/module"
	hello "github.com/aakashRajur/star-wars/internal/wiki/api/hello/http"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
)

func GetResources(resourceGroup di.ResourcesCompiler) []http.Resource {
	return resourceGroup.Resources
}

var ResourceModule = fx.Provide(
	module.StatsResource,
	hello.Resource,
	GetResources,
)
