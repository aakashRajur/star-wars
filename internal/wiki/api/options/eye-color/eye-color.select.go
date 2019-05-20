package eye_color

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
where constant_type = 'EYE_COLOR'
order by id;
`
)

func QueryEyeColor(storage types.Storage, tracker types.TimeTracker) ([]map[string]interface{}, error) {
	defer tracker(time.Now())
	return storage.GetArray(
		`eye_color`,
		types.Query{
			QueryString: Query,
			Args:        []interface{}{},
		},
	)
}
