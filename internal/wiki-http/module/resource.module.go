package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/api/hello"
	"github.com/aakashRajur/star-wars/internal/common/api/stats"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/character"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/characters"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/film"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/films"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/planet"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/planets"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/specie"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/species"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/vehicle"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/vehicles"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
)

func GetResources(resourceGroup di.HttpResourcesCompiler) []http.Resource {
	return resourceGroup.Resources
}

var ResourceModule = fx.Provide(
	stats.Resource,
	hello.Resource,
	planets.Resource,
	planet.Resource,
	species.Resource,
	specie.Resource,
	vehicles.Resource,
	vehicle.Resource,
	characters.Resource,
	character.Resource,
	films.Resource,
	film.Resource,
	GetResources,
)
