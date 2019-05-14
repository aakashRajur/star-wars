package di

import (
	"github.com/aakashRajur/star-wars/pkg/service"
	"go.uber.org/fx"
)

type ResourceSubscriptionCompiler struct {
	fx.In
	Subscriptions []service.Subscription `group:"resource_subscriptions"`
}

type ResourceSubscriptionProvider struct {
	fx.Out
	Subscription *service.Subscription `group:"resource_subscriptions"`
}
