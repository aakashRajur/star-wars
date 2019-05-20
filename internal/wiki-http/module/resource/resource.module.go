package resource

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/wiki-http/api/options/climate"
	eye_color "github.com/aakashRajur/star-wars/internal/wiki-http/api/options/eye-color"
	hair_color "github.com/aakashRajur/star-wars/internal/wiki-http/api/options/hair-color"
	skin_color "github.com/aakashRajur/star-wars/internal/wiki-http/api/options/skin-color"
	"github.com/aakashRajur/star-wars/internal/wiki-http/api/options/terrain"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"

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
	"github.com/aakashRajur/star-wars/pkg/http"
)

func GetResources(resourceGroup http_resource.HttpResourcesCompiler) []http.Resource {
	return resourceGroup.Resources
}

var Module = fx.Provide(
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
	climate.Resource,
	eye_color.Resource,
	hair_color.Resource,
	skin_color.Resource,
	terrain.Resource,
	GetResources,
)
