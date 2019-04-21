package redis_pg

import (
	"fmt"
	"strings"

	"github.com/aakashRajur/star-wars/pkg/pg"
	"github.com/aakashRajur/star-wars/pkg/redis"
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/util"
)

const (
	CacheKeySeparator = `:`
	CacheKey          = `%s:%s`
)

func generateRedisKey(key string, args []interface{}) (*string, error) {
	casted, err := util.MapStringFromInterfaces(args)
	if err != nil {
		return nil, err
	}
	coalesced := fmt.Sprintf(
		CacheKey,
		key,
		strings.Join(
			casted,
			CacheKeySeparator,
		),
	)
	return &coalesced, nil
}

type RedisPg struct {
	redis           *redis.Redis
	pg              *pg.Pg
	updateListeners map[string][]types.UpdateListener
	logger          types.Logger
}

func (instance *RedisPg) GetLogger() types.Logger {
	return instance.logger
}

func (instance *RedisPg) Close() error {
	err1 := instance.redis.Close()
	err2 := instance.pg.Close()

	if err1 != nil {
		return err1
	}
	return err2
}

func (instance *RedisPg) GetDatabase() types.Database {
	return instance.pg
}

func (instance *RedisPg) GetCache() types.Cache {
	return instance.redis
}

func (instance *RedisPg) GetObject(key string, query types.Query) (map[string]interface{}, error) {
	redisKey, err := generateRedisKey(key, query.Args)
	if err != nil {
		return nil, err
	}
	data, err := instance.redis.GetObject(*redisKey)
	if err == nil {
		return data, nil
	}
	if !redis.IsNil(err) {
		return nil, err
	}

	data, err = instance.pg.GetObject(query)
	if err != nil {
		return nil, err
	}

	err = instance.redis.SetObject(*redisKey, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (instance *RedisPg) GetArray(key string, query types.Query) ([]map[string]interface{}, error) {
	redisKey, err := generateRedisKey(key, query.Args)
	if err != nil {
		return nil, err
	}
	data, err := instance.redis.GetArray(*redisKey)
	if err == nil {
		return data, nil
	}
	if !redis.IsNil(err) {
		return nil, err
	}

	data, err = instance.pg.GetArray(query)
	if err != nil {
		return nil, err
	}

	err = instance.redis.SetArray(*redisKey, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (instance *RedisPg) GetPaginatedArray(key string, query types.Query, pagination types.Pagination, recordIdKey string) ([]map[string]interface{}, *types.Pagination, error) {
	data, page, err := instance.redis.GetPaginatedArray(key, pagination)
	if err == nil {
		return data, page, nil
	}
	if !redis.IsNil(err) {
		return nil, nil, err
	}

	data, page, err = instance.pg.GetPaginatedArray(query, pagination, recordIdKey)
	if err != nil {
		return nil, nil, err
	}

	err = instance.redis.SetPaginatedArray(key, pagination, data, *page)
	if err != nil {
		return nil, nil, err
	}

	return data, page, nil
}

func (instance *RedisPg) Set(queries []types.Query) error {
	return instance.pg.Set(queries)
}

func (instance *RedisPg) Listen(listener types.UpdateListener, channels ...string) error {
	for _, channel := range channels {
		listeners, ok := instance.updateListeners[channel]
		if !ok {
			listeners = []types.UpdateListener{listener}
			instance.updateListeners[channel] = listeners
			if channel != `*` {
				err := instance.pg.Listen(channel)
				if err != nil {
					return err
				}
			}
		} else {
			instance.updateListeners[channel] = append(listeners, listener)
		}
	}
	return nil
}

func (instance *RedisPg) OnNotification(notification types.Notification) {
	listeners := instance.updateListeners[notification.GetChannel()]
	global, ok := instance.updateListeners[`*`]
	if ok {
		listeners = append(global, listeners...)
	}
	for _, each := range listeners {
		if listeners != nil {
			each(instance, notification)
		}
	}
}

func (instance *RedisPg) GenerateCacheKey(key string, args ...interface{}) (string, error) {
	casted, err := util.MapStringFromInterfaces(args)
	if err != nil {
		return ``, err
	}
	coalesced := fmt.Sprintf(
		CacheKey,
		key,
		strings.Join(
			casted,
			CacheKeySeparator,
		),
	)
	return coalesced, nil
}
