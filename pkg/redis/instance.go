package redis

import (
	"github.com/go-redis/redis"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func NewInstance(uri string, logger types.Logger) (*Redis, error) {
	options, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	instance := Redis{redis.NewClient(options)}
	instance.WrapProcess(NewLogger(logger))

	return &instance, nil
}
