package consul

import (
	"context"
	"time"

	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func NewInstance(config Config, logger types.Logger) *Consul {
	ctx, cancel := context.WithCancel(context.Background())
	consul := &Consul{
		logger:   logger,
		config:   config,
		client:   http.NewTransport(10 * time.Second),
		ctx:      ctx,
		cancel:   cancel,
		services: make(map[string][]string),
		ready:    make(chan bool),
	}
	go consul.watchServices()
	return consul
}
