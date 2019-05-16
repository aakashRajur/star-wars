package vehicle

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	Query = `select * from vehicles where id = $1;`
)

func QuerySelectVehicle(storage types.Storage, tracker types.TimeTracker, cacheKey string, id int) (map[string]interface{}, error) {
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
