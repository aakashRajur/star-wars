package redis_pg

import (
	"encoding/json"

	"github.com/aakashRajur/star-wars/pkg/types"
)

type CacheStrategy struct {
	Channel  string
	CacheKey string
}

func CacheBuster(strategies []CacheStrategy) types.UpdateListener {
	return func(storage types.Storage, notification types.Notification) {
		channel := notification.GetChannel()
		var strategy *CacheStrategy = nil
		for _, each := range strategies {
			if each.Channel == channel {
				strategy = &each
				break
			}
		}
		if strategy == nil {
			return
		}

		cacheKey := strategy.CacheKey

		logger := storage.GetLogger()

		logger.InfoFields(
			map[string]interface{}{
				`channel`:   cacheKey,
				`payload`:   notification.GetPayload(),
				`timestamp`: notification.GetTimestamp(),
			},
		)

		ids := make([]int64, 0)
		err := json.Unmarshal([]byte(notification.GetPayload()), &ids)
		if err != nil {
			return
		}

		cache := storage.GetCache()
		objectKeys := make([]string, 0)
		paginationKeys := make([]string, 0)

		for _, each := range ids {
			generated, err := storage.GenerateCacheKey(cacheKey, each)
			if err == nil {
				objectKeys = append(
					objectKeys,
					generated,
				)
			}
			compiled, err := cache.GeneratePaginationCacheKeysForId(cacheKey, each)
			if err == nil {
				paginationKeys = append(paginationKeys, compiled...)
			}
		}

		logger.InfoFields(
			map[string]interface{}{
				`objectkey`:      objectKeys,
				`paginationkeys`: paginationKeys,
			},
		)

		err = cache.Delete(objectKeys...)
		if err != nil {
			return
		}

		err = cache.DeletePagination(paginationKeys...)
		if err != nil {
			return
		}
	}
}
