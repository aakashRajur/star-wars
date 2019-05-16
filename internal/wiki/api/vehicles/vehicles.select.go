package vehicles

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	Query        = `select id, name, model from vehicles order by id;`
	PageRecordId = `id`
)

func QuerySelectVehicles(storage types.Storage, tracker types.TimeTracker, cacheKey string, pagination types.Pagination) ([]map[string]interface{}, *types.Pagination, error) {
	defer tracker(time.Now())

	data, page, err := storage.GetPaginatedArray(
		cacheKey,
		types.Query{
			QueryString: Query,
			Args:        []interface{}{},
		},
		pagination,
		PageRecordId,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}
