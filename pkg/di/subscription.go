package di

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/kafka"
)

type SubscriptionCompiler struct {
	fx.In
	Subscriptions []*kafka.Subscription `group:"subscriptions"`
}

type SubscriptionProvider struct {
	fx.Out
	Subscription *kafka.Subscription `group:"subscriptions"`
}
