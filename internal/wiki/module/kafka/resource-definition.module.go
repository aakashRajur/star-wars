package kafka

import (
	"go.uber.org/fx"

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
	"github.com/aakashRajur/star-wars/pkg/resource-definition"
)

func GetResourceDefinitions(definitionCompiler di.ResourceDefinitionCompiler) []resource_definition.ResourceDefinition {
	return definitionCompiler.ResourceDefinitions
}

var ResourceDefinitionModule = fx.Provide(
	character.ProvideResourceDefinitionGet,
	character.ProvideResourceDefinitionPatch,
	characters.ProvideResourceDefinitionGet,
	film.ProvideResourceDefinitionGet,
	film.ProvideResourceDefinitionPatch,
	films.ProvideResourceDefinitionGet,
	planet.ProvideResourceDefinitionGet,
	planet.ProvideResourceDefinitionPatch,
	planets.ProvideResourceDefinitionGet,
	specie.ProvideResourceDefinitionGet,
	specie.ProvideResourceDefinitionPatch,
	species.ProvideResourceDefinitionGet,
	vehicle.ProvideResourceDefinitionGet,
	vehicle.ProvideResourceDefinitionPatch,
	vehicles.ProvideResourceDefinitionGet,
	GetResourceDefinitions,
)
