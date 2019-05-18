package subscription

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/kafka"
)

func GetSubscriptions() []*kafka.Subscription {
	return []*kafka.Subscription{}
}

var Module = fx.Provide(GetSubscriptions)
