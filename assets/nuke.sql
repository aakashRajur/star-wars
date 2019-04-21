drop trigger if exists films_insert_created on films;
  drop trigger if exists films_update_updated on films;
  drop trigger if exists films_notify_updated on films;
  drop trigger if exists films_notify_deleted on films;

  drop trigger if exists characters_insert_created on characters;
  drop trigger if exists characters_update_updated on characters;
  drop trigger if exists characters_notify_updated on films;
  drop trigger if exists characters_notify_deleted on films;

  drop trigger if exists vehicles_insert_created on vehicles;
  drop trigger if exists vehicles_update_updated on vehicles;
  drop trigger if exists vehicles_notify_updated on films;
  drop trigger if exists vehicles_notify_deleted on films;

  drop trigger if exists species_insert_created on species;
  drop trigger if exists species_update_updated on species;
  drop trigger if exists species_notify_updated on films;
  drop trigger if exists species_notify_deleted on films;

  drop trigger if exists planets_insert_created on planets;
  drop trigger if exists planets_update_updated on planets;
  drop trigger if exists planets_notify_updated on films;
  drop trigger if exists planets_notify_deleted on films;

  drop trigger if exists normalized_relations_insert_created on normalized_relations;
  drop trigger if exists normalized_relations_update_updated on normalized_relations;
  drop trigger if exists normalized_relations_verify_insert_update on normalized_relations;

  drop trigger if exists constants_insert_created on constants;
  drop trigger if exists constants_update_updated on constants;

  drop table if exists films;
  drop table if exists characters;
  drop table if exists vehicles;
  drop table if exists species;
  drop table if exists planets;
  drop table if exists normalized_relations;
  drop table if exists constants;

  drop function if exists create_insert_created;
  drop function if exists create_update_updated;
  drop function if exists create_verify_normalized_relations;
  drop function if exists notify_updated;
  drop function if exists notify_deleted;

drop database if exists star_wars;
