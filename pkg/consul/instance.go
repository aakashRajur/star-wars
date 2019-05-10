package consul

import (
	"github.com/aakashRajur/star-wars/pkg/types"
)

func NewInstance(config Config, logger types.Logger) *Consul {
	consul := Consul{
		Config: config,
		Logger: logger,
	}
	return &consul
}
