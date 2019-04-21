package species

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	SpeciesQuery        = `select id, name, classification, spoken_language from species order by id;`
	SpeciesPageRecordId = `id`
)

func QuerySelectSpecies(storage types.Storage, tracker types.TimeTracker, cacheKey string, pagination types.Pagination) ([]map[string]interface{}, *types.Pagination, error) {
	defer tracker(time.Now())
	data, page, err := storage.GetPaginatedArray(
		cacheKey,
		types.Query{
			QueryString: SpeciesQuery,
			Args:        []interface{}{},
		},
		pagination,
		SpeciesPageRecordId,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}
