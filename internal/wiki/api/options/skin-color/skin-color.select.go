package skin_color

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
where constant_type = 'SKIN_COLOR'
order by id;
`
)

func QuerySkinColor(storage types.Storage, tracker types.TimeTracker) ([]map[string]interface{}, error) {
	defer tracker(time.Now())
	return storage.GetArray(
		`skin_color`,
		types.Query{
			QueryString: Query,
			Args:        []interface{}{},
		},
	)
}
