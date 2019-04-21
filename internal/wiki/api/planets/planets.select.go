package planets

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	PlanetsQuery        = `select id, name from planets order by id;`
	PlanetsPageRecordId = `id`
)

func QuerySelectPlanets(storage types.Storage, tracker types.TimeTracker, cacheKey string, pagination types.Pagination) ([]map[string]interface{}, *types.Pagination, error) {
	defer tracker(time.Now())
	data, page, err := storage.GetPaginatedArray(
		cacheKey,
		types.Query{
			QueryString: PlanetsQuery,
			Args:        []interface{}{},
		},
		pagination,
		PlanetsPageRecordId,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}
