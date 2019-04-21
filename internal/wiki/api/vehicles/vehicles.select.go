package vehicles

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	VehiclesQuery        = `select id, name, model from vehicles order by id;`
	VehiclesPageRecordId = `id`
)

func QuerySelectVehicles(storage types.Storage, tracker types.TimeTracker, cacheKey string, pagination types.Pagination) ([]map[string]interface{}, *types.Pagination, error) {
	defer tracker(time.Now())

	data, page, err := storage.GetPaginatedArray(
		cacheKey,
		types.Query{
			QueryString: VehiclesQuery,
			Args:        []interface{}{},
		},
		pagination,
		VehiclesPageRecordId,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}
