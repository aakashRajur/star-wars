package consul

import (
	"context"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func NewInstance(config Config, logger types.Logger) *Consul {
	ctx, cancel := context.WithCancel(context.Background())
	consul := Consul{
		config:       config,
		logger:       logger,
		ctx:          ctx,
		cancel:       cancel,
		observations: make(map[string][]service.Subscription),
	}
	return &consul
}
