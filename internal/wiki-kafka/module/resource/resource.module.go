package resource

import (
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/api/hello"
	"github.com/aakashRajur/star-wars/internal/common/api/stats"
	"github.com/aakashRajur/star-wars/pkg/http"
)

func GetResources(resourceGroup http_resource.HttpResourcesCompiler) []http.Resource {
	return resourceGroup.Resources
}

var Module = fx.Provide(
	stats.Resource,
	hello.Resource,
	GetResources,
)
