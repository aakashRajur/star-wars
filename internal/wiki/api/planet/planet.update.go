package planet

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func QueryUpdatePlanet(storage types.Storage, tracker types.TimeTracker, id int, fields map[string]interface{}) error {
	defer tracker(time.Now())

	db := storage.GetDatabase()

	updateQuery := db.GenerateUpdateQuery(
		`planets`,
		fields,
		[]types.Constraint{
			{
				Field:    `id`,
				Relation: `=`,
				Value:    id,
			},
		},
	)

	return storage.Set([]types.Query{updateQuery})
}
