package characters

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	Query        = `select id, name, gender, birth_year from characters order by id;`
	PageRecordId = `id`
)

func QuerySelectCharacters(storage types.Storage, tracker types.TimeTracker, cacheKey string, pagination types.Pagination) ([]map[string]interface{}, *types.Pagination, error) {
	defer tracker(time.Now())
	return storage.GetPaginatedArray(
		cacheKey,
		types.Query{
			QueryString: Query,
			Args:        []interface{}{},
		},
		pagination,
		PageRecordId,
	)
}
