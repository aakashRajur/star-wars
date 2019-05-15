package consul

import (
	"context"
	"net/http"
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func NewInstance(config Config, logger types.Logger) *Consul {
	ctx, cancel := context.WithCancel(context.Background())
	consul := &Consul{
		logger: logger,
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 10,
				MaxIdleConns:        10,
				IdleConnTimeout:     1 * time.Second,
			},
		},
		ctx:      ctx,
		cancel:   cancel,
		services: make(map[string][]string),
		ready:    make(chan bool),
	}
	go consul.watchServices()
	return consul
}
