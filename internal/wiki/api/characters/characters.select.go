package characters

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	CharactersQuery       = `select id, name, gender, birth_year from characters order by id;`
	CharacterPageRecordId = `id`
)

func QuerySelectCharacters(storage types.Storage, tracker types.TimeTracker, cacheKey string, pagination types.Pagination) ([]map[string]interface{}, *types.Pagination, error) {
	defer tracker(time.Now())
	data, page, err := storage.GetPaginatedArray(
		cacheKey,
		types.Query{
			QueryString: CharactersQuery,
			Args:        []interface{}{},
		},
		pagination,
		CharacterPageRecordId,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}
