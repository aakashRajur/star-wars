package module

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/character"
	"github.com/aakashRajur/star-wars/internal/wiki/api/characters"
	"github.com/aakashRajur/star-wars/internal/wiki/api/film"
	"github.com/aakashRajur/star-wars/internal/wiki/api/films"
	"github.com/aakashRajur/star-wars/internal/wiki/api/planet"
	"github.com/aakashRajur/star-wars/internal/wiki/api/planets"
	"github.com/aakashRajur/star-wars/internal/wiki/api/specie"
	"github.com/aakashRajur/star-wars/internal/wiki/api/species"
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicle"
	"github.com/aakashRajur/star-wars/internal/wiki/api/vehicles"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/service"
	"go.uber.org/fx"
)

func GetResources(resourceGroup di.ServiceResourceCompiler) []service.Resource {
	return resourceGroup.Resources
}

var ResourceModule = fx.Provide(
	character.Get,
	character.Patch,
	characters.Get,
	film.Get,
	film.Patch,
	films.Get,
	planet.Get,
	planet.Patch,
	planets.Get,
	specie.Get,
	specie.Patch,
	species.Get,
	vehicle.Get,
	vehicle.Patch,
	vehicles.Get,
	GetResources,
)
