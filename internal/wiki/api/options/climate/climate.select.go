package climate

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	Query = `
select
    id,
    constant_value,
    constant_type
from constants
where constant_type = 'CLIMATE'
order by id;
`
)

func QueryClimate(storage types.Storage, tracker types.TimeTracker) ([]map[string]interface{}, error) {
	defer tracker(time.Now())
	return storage.GetArray(
		`climate`,
		types.Query{
			QueryString: Query,
			Args:        []interface{}{},
		},
	)
}
