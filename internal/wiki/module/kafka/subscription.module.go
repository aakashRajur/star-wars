package kafka

import (
	"go.uber.org/fx"

	character "github.com/aakashRajur/star-wars/internal/wiki/api/character/kafka"
	characters "github.com/aakashRajur/star-wars/internal/wiki/api/characters/kafka"
	film "github.com/aakashRajur/star-wars/internal/wiki/api/film/kafka"
	films "github.com/aakashRajur/star-wars/internal/wiki/api/films/kafka"
	planet "github.com/aakashRajur/star-wars/internal/wiki/api/planet/kafka"
	planets "github.com/aakashRajur/star-wars/internal/wiki/api/planets/kafka"
	specie "github.com/aakashRajur/star-wars/internal/wiki/api/specie/kafka"
	species "github.com/aakashRajur/star-wars/internal/wiki/api/species/kafka"
	vehicle "github.com/aakashRajur/star-wars/internal/wiki/api/vehicle/kafka"
	vehicles "github.com/aakashRajur/star-wars/internal/wiki/api/vehicles/kafka"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/kafka"
)

func GetSubscriptions(subscriptionGroup di.SubscriptionCompiler) []*kafka.Subscription {
	return subscriptionGroup.Subscriptions
}

var SubscriptionModule = fx.Provide(
	character.GetCharacter,
	character.PatchCharacter,
	characters.GetCharacters,
	film.GetFilm,
	film.PatchFilm,
	films.GetFilms,
	planet.GetPlanet,
	planet.PatchPlanet,
	planets.GetPlanets,
	specie.GetSpecie,
	specie.PatchSpecie,
	species.GetSpecies,
	vehicle.GetVehicle,
	vehicle.PatchVehicle,
	vehicles.GetVehicles,
	GetSubscriptions,
)
