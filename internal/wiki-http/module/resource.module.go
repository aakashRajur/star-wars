package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/common/api/hello"
	"github.com/aakashRajur/star-wars/internal/common/api/stats"
	"github.com/aakashRajur/star-wars/internal/wiki/api/character/http"
	"github.com/aakashRajur/star-wars/internal/wiki/api/characters/http"
	"github.com/aakashRajur/star-wars/internal/wiki/api/film/http"
	"github.com/aakashRajur/star-wars/internal/wiki/api/films/http"
	"github.com/aakashRajur/star-wars/internal/wiki/api/planet/http"
	"github.com/aakashRajur/star-wars/internal/wiki/api/planets/http"
	"github.com/aakashRajur/star-wars/internal/wiki/api/specie/http"
	"github.com/aakashRajur/star-wars/internal/wiki/api/species/http"
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicle/http"
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicles/http"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
)

func GetResources(resourceGroup di.ResourcesCompiler) []http.Resource {
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
