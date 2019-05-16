package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/api-gateway/api/character"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/film"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/planet"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/specie"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/vehicle"
	"github.com/aakashRajur/star-wars/internal/common/api/hello"
	"github.com/aakashRajur/star-wars/internal/common/api/stats"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
)

func GetHttpResources(resourceGroup di.HttpResourcesCompiler) []http.Resource {
	return resourceGroup.Resources
}

var ResourceModule = fx.Provide(
	stats.Resource,
	hello.Resource,
	character.HttpResource,
	film.HttpResource,
	planet.HttpResource,
	specie.HttpResource,
	vehicle.HttpResource,
	GetHttpResources,
)
