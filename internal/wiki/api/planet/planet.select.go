package planet

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	PlanetQuery = `select * from planets where id = $1;`
)

func QuerySelectPlanet(storage types.Storage, tracker types.TimeTracker, cacheKey string, id int) (map[string]interface{}, error) {
	defer tracker(time.Now())

	data, err := storage.GetObject(
		cacheKey,
		types.Query{
			QueryString: PlanetQuery,
			Args:        []interface{}{id},
		},
	)
	if err != nil {
		return nil, err
	}
	return data, nil
}
