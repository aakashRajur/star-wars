package resource

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/api-gateway/api/async"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"

	"github.com/aakashRajur/star-wars/internal/api-gateway/api/character"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/characters"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/film"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/films"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/planet"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/planets"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/specie"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/species"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/vehicle"
	"github.com/aakashRajur/star-wars/internal/api-gateway/api/vehicles"
	"github.com/aakashRajur/star-wars/internal/common/api/hello"
	"github.com/aakashRajur/star-wars/internal/common/api/stats"
	"github.com/aakashRajur/star-wars/pkg/http"
)

func GetHttpResources(resourceGroup http_resource.HttpResourcesCompiler) []http.Resource {
	return resourceGroup.Resources
}

var Module = fx.Provide(
	stats.Resource,
	hello.Resource,
	character.HttpResource,
	character.KafkaResource,
	characters.HttpResource,
	characters.KafkaResource,
	film.HttpResource,
	film.KafkaResource,
	films.HttpResource,
	films.KafkaResource,
	planet.HttpResource,
	planet.KafkaResource,
	planets.HttpResource,
	planets.KafkaResource,
	specie.HttpResource,
	specie.KafkaResource,
	species.HttpResource,
	species.KafkaResource,
	vehicle.HttpResource,
	vehicle.KafkaResource,
	vehicles.HttpResource,
	vehicles.KafkaResource,
	async.HttpResource,
	GetHttpResources,
)
