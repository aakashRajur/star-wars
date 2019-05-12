package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/wiki-kafka/api/character"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/api/characters"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/api/film"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/api/films"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/api/planet"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/api/planets"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/api/specie"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/api/species"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/api/vehicle"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/api/vehicles"
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
