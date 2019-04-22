package kafka

import (
	"go.uber.org/fx"

	character "github.com/aakashRajur/star-wars/internal/wiki/api/character/kafka"
	film "github.com/aakashRajur/star-wars/internal/wiki/api/film/kafka"
	planet "github.com/aakashRajur/star-wars/internal/wiki/api/planet/kafka"
	specie "github.com/aakashRajur/star-wars/internal/wiki/api/specie/kafka"
	vehicle "github.com/aakashRajur/star-wars/internal/wiki/api/vehicle/kafka"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/kafka"
)

func GetSubscriptions(subscriptionGroup di.SubscriptionCompiler) []*kafka.Subscription {
	return subscriptionGroup.Subscriptions
}

var SubscriptionModule = fx.Provide(
	character.GetCharacter,
	character.PatchCharacter,
	film.GetFilm,
	film.PatchFilm,
	planet.GetPlanet,
	planet.PatchPlanet,
	specie.GetSpecie,
	specie.PatchSpecie,
	vehicle.GetVehicle,
	vehicle.PatchVehicle,
	GetSubscriptions,
)
