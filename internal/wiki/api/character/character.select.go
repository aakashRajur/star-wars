package character

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	CharacterQuery = `
with
  vehicle_ids   as (
    select
      unnest(characters.vehicles) as id
    from characters
    where id = $1
  ),
  vehicles_json as (
    select
      $1                  as id,
      json_agg(_vehicles) as vehicles
    from (
     select
       vehicles.id,
       vehicles.name
     from vehicle_ids
            inner join vehicles on vehicle_ids.id = vehicles.id
     order by vehicles.id
   ) as _vehicles
  )
select
  characters.id,
  characters.name,
  characters.mass,
  characters.hair_color,
  characters.skin_color,
  characters.eye_color,
  characters.birth_year,
  characters.gender,
  json_build_object(planets.id, planets.name) as home_world,
  vehicles_json.vehicles,
  json_build_object(species.id, species.name) as species,
  characters.description
from characters
       inner join vehicles_json on vehicles_json.id = characters.id
       inner join planets on planets.id = characters.id
       inner join species on characters.species = species.id
where characters.id = $1;
`
)

func QuerySelectCharacter(storage types.Storage, tracker types.TimeTracker, cacheKey string, id int) (map[string]interface{}, error) {
	defer tracker(time.Now())

	data, err := storage.GetObject(
		cacheKey,
		types.Query{
			QueryString: CharacterQuery,
			Args:        []interface{}{id},
		},
	)
	if err != nil {
		return nil, err
	}
	return data, nil
}
