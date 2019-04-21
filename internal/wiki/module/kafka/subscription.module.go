package kafka

import (
	"go.uber.org/fx"

	character "github.com/aakashRajur/star-wars/internal/wiki/api/character/kafka"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/kafka"
)

func GetSubscriptions(subscriptionGroup di.SubscriptionCompiler) []*kafka.Subscription {
	return subscriptionGroup.Subscriptions
}

var SubscriptionModule = fx.Provide(
	character.GetCharacter,
	character.PatchCharacter,
	GetSubscriptions,
)
