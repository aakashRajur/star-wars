package film

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	FilmQuery = `
with
  planet_ids      as (
    select
      unnest(films.planets) as id
    from films
    where id = $1
  ),
  planets_json    as (
    select
      $1                  as id,
      json_agg(_planets)  as planets
    from (
           select
             planets.id,
             planets.name
           from planet_ids
                  inner join planets on planet_ids.id = planets.id
           order by planets.id
         ) as _planets
  ),
  vehicles_ids    as (
    select
      unnest(films.vehicles) as id
    from films
    where id = $1
  ),
  vehicles_json   as (
    select
      $1                   as id,
      json_agg(_vehicles)  as vehicles
    from (
           select
             vehicles.id,
             vehicles.name
           from vehicles_ids
                  inner join vehicles on vehicles_ids.id = vehicles.id
           order by vehicles.id
         ) as _vehicles
  ),
  star_ship_ids   as (
    select
      unnest(films.star_ships) as id
    from films
    where id = $1
  ),
  star_ships_json as (
    select
      $1                     as id,
      json_agg(_star_ships)  as star_ships
    from (
           select
             vehicles.id,
             vehicles.name
           from star_ship_ids
                  inner join vehicles on star_ship_ids.id = vehicles.id
           order by vehicles.id
         ) as _star_ships
  ),
  specie_ids      as (
    select
      unnest(films.species) as id
    from films
    where id = $1
  ),
  species_json    as (
    select
      $1                  as id,
      json_agg(_species)  as species
    from (
           select
             species.id,
             species.name
           from specie_ids
                  inner join species on specie_ids.id = species.id
           order by species.id
         ) as _species
  ),
  character_ids   as (
    select
      unnest(films.characters) as id
    from films
    where id = $1
  ),
  characters_json as (
    select
      $1                     as id,
      json_agg(_characters)  as characters
    from (
           select
             characters.id,
             characters.name
           from character_ids
                  inner join characters on character_ids.id = characters.id
           order by characters.id
         ) as _characters
  )

select
  films.id,
  films.title,
  films.episode,
  films.opening_crawl,
  films.director,
  films.producer,
  films.release_date,
  films.description,
  planets_json.planets       as planets,
  vehicles_json.vehicles     as vehicles,
  star_ships_json.star_ships as star_ships,
  species_json.species       as species,
  characters_json.characters as characters
from films
       inner join planets_json on planets_json.id = films.id
       inner join vehicles_json on vehicles_json.id = films.id
       inner join star_ships_json on star_ships_json.id = films.id
       inner join species_json on species_json.id = films.id
       inner join characters_json on characters_json.id = films.id
where films.id = $1;
`
)

func QuerySelectFilm(storage types.Storage, tracker types.TimeTracker, cacheKey string, id int) (map[string]interface{}, error) {
	defer tracker(time.Now())

	data, err := storage.GetObject(
		cacheKey,
		types.Query{
			QueryString: FilmQuery,
			Args:        []interface{}{id},
		},
	)
	if err != nil {
		return nil, err
	}
	return data, nil
}
