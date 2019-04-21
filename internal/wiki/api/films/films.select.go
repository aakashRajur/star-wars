package films

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	FilmsQuery        = `select id, title, episode, director, release_date from films order by id;`
	FilmsPageRecordId = `id`
)

func QuerySelectFilms(storage types.Storage, tracker types.TimeTracker, cacheKey string, pagination types.Pagination) ([]map[string]interface{}, *types.Pagination, error) {
	defer tracker(time.Now())
	data, page, err := storage.GetPaginatedArray(
		cacheKey,
		types.Query{
			QueryString: FilmsQuery,
			Args:        []interface{}{},
		},
		pagination,
		FilmsPageRecordId,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}
