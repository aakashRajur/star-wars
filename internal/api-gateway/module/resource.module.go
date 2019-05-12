package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/api/hello"
	"github.com/aakashRajur/star-wars/internal/common/api/stats"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
)

func GetResources(resourceGroup di.HttpResourcesCompiler) []http.Resource {
	return resourceGroup.Resources
}

var ResourceModule = fx.Provide(
	stats.Resource,
	hello.Resource,
	GetResources,
)
