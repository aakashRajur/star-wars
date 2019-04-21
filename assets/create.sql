create or replace function create_verify_normalized_relations()
  returns trigger as
$BODY$
declare
  primary_index   int;
  secondary_index int;
begin
  if new.primary_table is not null and not exists(
      select
        1
      from pg_catalog.pg_tables
      where tablename = new.primary_table and schemaname = 'public'
    )
  then
    raise notice 'table (%) does not exist', new.primary_table;
  end if;

  if new.secondary_table is not null and not exists(
      select
        1
      from pg_catalog.pg_tables
      where tablename = new.secondary_table and schemaname = 'public'
    )
  then
    raise notice 'table (%) does not exist', new.primary_table;
  end if;

  if new.primary_index is not null then
    execute format(
        'select id from %I where id = %s',
        new.primary_table,
        new.primary_index
      ) into primary_index;

    if primary_index is null
    then
      raise notice '(%I) does not exist in table (%)', new.primary_index, new.primary_table;
    end if;
  end if;

  if new.secondary_index is not null
  then
    execute format(
        'select id from %I where id = %s',
        new.secondary_table,
        new.secondary_index
      ) into secondary_index;

    if secondary_index is null
    then
      raise notice '(%I) does not exist in table (%)', new.secondary_index, new.secondary_table;
    end if;
  end if;

  return new;
end;
$BODY$
  language plpgsql;

create or replace function create_insert_created()
  returns trigger as
$BODY$
declare
  _now timestamp := current_timestamp;
begin
  new.updated = _now;
  new.created = _now;
  return new;
end;
$BODY$
  language plpgsql;

create or replace function create_update_updated()
  returns trigger as
$BODY$
begin
  new.updated = current_timestamp;
  return new;
end;
$BODY$
  language plpgsql;

create or replace function notify_updated()
  returns trigger as
$BODY$
declare
  channel text := TG_ARGV [ 0];
begin
  perform
  (
    with
      updated as (
        select
          id
        from new_table
        order by id
        ),
      stringified as (select json_agg(id)::text as updated from updated)
      select pg_notify(channel, (select updated from stringified))
  );
  return null;
end;
$BODY$
  language plpgsql;

create or replace function notify_deleted()
  returns trigger as
$BODY$
declare
  channel text := TG_ARGV [ 0];
begin
  perform
  (
    with
      deleted as (
        select
          id
        from old_table
        order by id
        ),
      stringified as (select json_agg(id)::text as deleted from deleted)
      select pg_notify(channel, (select deleted from stringified))
  );
  return null;
end;
$BODY$
  language plpgsql;

create table if not exists constants (
                                       id             bigserial primary key,
                                       constant_type  text      not null,
                                       constant_value text      not null,
                                       updated        timestamp not null,
                                       created        timestamp not null,
                                       constraint constants_unique unique (constant_type, constant_value)
);

create trigger constants_insert_created
  before insert
  on constants
  for each row
execute procedure create_insert_created();

create trigger constants_update_updated
  before update
  on constants
  for each row
execute procedure create_update_updated();

create table if not exists planets (
                                     id              bigserial primary key,
                                     name            text not null,
                                     rotation_period float,
                                     orbital_period  float,
                                     diameter        float,
                                     climate         text[] default array []::text[],
                                     terrain         text[] default array []::text[],
                                     gravity         float,
                                     surface_water   float,
                                     population      bigint,
                                     description     text,
                                     created         timestamp,
                                     updated         timestamp
);

create trigger planets_insert_created
  before insert
  on planets
  for each row
execute procedure create_insert_created();

create trigger planets_update_updated
  before update
  on planets
  for each row
execute procedure create_update_updated();

create trigger planets_notify_updated
  after update
  on planets
  referencing new table as new_table
  for each statement
execute procedure notify_updated('planets');

create trigger planets_notify_deleted
  after delete
  on planets
  referencing old table as old_table
  for each statement
execute procedure notify_deleted('planets');

create table if not exists species (
                                     id               bigserial primary key,
                                     name             text      not null,
                                     classification   text,
                                     average_height   float,
                                     skin_colors      text[] default array []::text[],
                                     hair_colors      text[] default array []::text[],
                                     eye_colors       text[] default array []::text[],
                                     average_lifespan float,
                                     home_world       bigint,
                                     spoken_language  text,
                                     description      text,
                                     created          timestamp not null,
                                     updated          timestamp not null,
                                     constraint species_home_world foreign key (home_world)
                                       references planets (id) match simple
                                       on update no action on delete no action
);

create trigger species_insert_created
  before insert
  on species
  for each row
execute procedure create_insert_created();

create trigger species_update_updated
  before update
  on species
  for each row
execute procedure create_update_updated();

create trigger species_notify_updated
  after update
  on species
  referencing new table as new_table
  for each statement
execute procedure notify_updated('species');

create trigger species_notify_deleted
  after delete
  on species
  referencing old table as old_table
  for each statement
execute procedure notify_deleted('species');

create table if not exists vehicles (
                                      id                    bigserial primary key,
                                      name                  text      not null,
                                      model                 text      not null,
                                      manufacturer          text,
                                      cost_in_credits       bigint,
                                      size                  float,
                                      max_atmospheric_speed float,
                                      crew                  int,
                                      passengers            int,
                                      cargo_capacity        bigint,
                                      consumables           interval,
                                      hyperdrive_rating     bytea,
                                      mglt                  int,
                                      starship_class        text,
                                      created               timestamp not null,
                                      updated               timestamp not null
);

create trigger vehicles_insert_created
  before insert
  on vehicles
  for each row
execute procedure create_insert_created();

create trigger vehicles_update_updated
  before update
  on vehicles
  for each row
execute procedure create_update_updated();

create trigger vehicles_notify_updated
  after update
  on vehicles
  referencing new table as new_table
  for each statement
execute procedure notify_updated('vehicles');

create trigger vehicles_notify_deleted
  after delete
  on vehicles
  referencing old table as old_table
  for each statement
execute procedure notify_deleted('vehicles');

create table if not exists characters (
                                        id          bigserial primary key,
                                        name        text      not null,
                                        height      int,
                                        mass        float,
                                        hair_color  bigint,
                                        skin_color  bigint,
                                        eye_color   bigint,
                                        birth_year  varchar(10),
                                        gender      varchar(10),
                                        home_world  bigint,
                                        species     bigint,
                                        vehicles    int[] default array []::int[],
                                        description text,
                                        created     timestamp not null,
                                        updated     timestamp not null,
                                        constraint character_home_world foreign key (home_world)
                                          references planets (id) match simple
                                          on update no action on delete no action,
                                        constraint character_species foreign key (species)
                                          references species (id) match simple
                                          on update no action on delete no action
);

create trigger characters_insert_created
  before insert
  on characters
  for each row
execute procedure create_insert_created();

create trigger characters_update_updated
  before update
  on characters
  for each row
execute procedure create_update_updated();

create trigger characters_notify_updated
  after update
  on characters
  referencing new table as new_table
  for each statement
execute procedure notify_updated('characters');

create trigger characters_notify_deleted
  after delete
  on characters
  referencing old table as old_table
  for each statement
execute procedure notify_deleted('characters');

create table if not exists films (
                                   id            bigserial primary key,
                                   title         text      not null,
                                   episode       int       not null,
                                   opening_crawl text      not null,
                                   director      text      not null,
                                   producer      text      not null,
                                   release_date  timestamp not null,
                                   characters    int[] default array []::int[],
                                   star_ships    int[] default array []::int[],
                                   vehicles      int[] default array []::int[],
                                   species       int[] default array []::int[],
                                   planets       int[] default array []::int[],
                                   description   text,
                                   created       timestamp not null,
                                   updated       timestamp not null
);

create trigger films_insert_created
  before insert
  on films
  for each row
execute procedure create_insert_created();

create trigger films_update_updated
  before update
  on films
  for each row
execute procedure create_update_updated();

create trigger films_notify_updated
  after update
  on films
  referencing new table as new_table
  for each statement
execute procedure notify_updated('films');

create trigger films_notify_deleted
  after delete
  on films
  referencing old table as old_table
  for each statement
execute procedure notify_deleted('films');
