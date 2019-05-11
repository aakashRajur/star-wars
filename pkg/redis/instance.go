package redis

import (
	"github.com/go-redis/redis"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func NewInstance(url Url, logger types.Logger) (*Redis, error) {
	options, err := redis.ParseURL(string(url))
	if err != nil {
		return nil, err
	}
	instance := Redis{redis.NewClient(options)}
	instance.WrapProcess(NewLogger(logger))

	return &instance, nil
}
