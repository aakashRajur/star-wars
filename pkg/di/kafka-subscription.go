package di

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/kafka"
)

type KafkaSubscriptionCompiler struct {
	fx.In
	Subscriptions []*kafka.Subscription `group:"kafka_subscriptions"`
}

type KafkaSubscriptionProvider struct {
	fx.Out
	Subscription *kafka.Subscription `group:"kafka_subscriptions"`
}
