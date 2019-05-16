package specie

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	Query = `select * from species where id = $1;`
)

func QuerySelectSpecie(storage types.Storage, tracker types.TimeTracker, cacheKey string, id int) (map[string]interface{}, error) {
	defer tracker(time.Now())

	data, err := storage.GetObject(
		cacheKey,
		types.Query{
			QueryString: Query,
			Args:        []interface{}{id},
		},
	)
	if err != nil {
		return nil, err
	}
	return data, nil
}
