--
-- PostgreSQL database dump
--

-- Dumped from database version 11.2
-- Dumped by pg_dump version 11.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE ONLY public.species DROP CONSTRAINT species_home_world;
ALTER TABLE ONLY public.characters DROP CONSTRAINT character_species;
ALTER TABLE ONLY public.characters DROP CONSTRAINT character_home_world;
DROP TRIGGER vehicles_update_updated ON public.vehicles;
DROP TRIGGER vehicles_notify_updated ON public.vehicles;
DROP TRIGGER vehicles_notify_deleted ON public.vehicles;
DROP TRIGGER vehicles_insert_created ON public.vehicles;
DROP TRIGGER species_update_updated ON public.species;
DROP TRIGGER species_notify_updated ON public.species;
DROP TRIGGER species_notify_deleted ON public.species;
DROP TRIGGER species_insert_created ON public.species;
DROP TRIGGER planets_update_updated ON public.planets;
DROP TRIGGER planets_notify_updated ON public.planets;
DROP TRIGGER planets_notify_deleted ON public.planets;
DROP TRIGGER planets_insert_created ON public.planets;
DROP TRIGGER films_update_updated ON public.films;
DROP TRIGGER films_notify_updated ON public.films;
DROP TRIGGER films_notify_deleted ON public.films;
DROP TRIGGER films_insert_created ON public.films;
DROP TRIGGER constants_update_updated ON public.constants;
DROP TRIGGER constants_insert_created ON public.constants;
DROP TRIGGER characters_update_updated ON public.characters;
DROP TRIGGER characters_notify_updated ON public.characters;
DROP TRIGGER characters_notify_deleted ON public.characters;
DROP TRIGGER characters_insert_created ON public.characters;
ALTER TABLE ONLY public.vehicles DROP CONSTRAINT vehicles_pkey;
ALTER TABLE ONLY public.species DROP CONSTRAINT species_pkey;
ALTER TABLE ONLY public.planets DROP CONSTRAINT planets_pkey;
ALTER TABLE ONLY public.films DROP CONSTRAINT films_pkey;
ALTER TABLE ONLY public.constants DROP CONSTRAINT constants_unique;
ALTER TABLE ONLY public.constants DROP CONSTRAINT constants_pkey;
ALTER TABLE ONLY public.characters DROP CONSTRAINT characters_pkey;
ALTER TABLE public.vehicles ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.species ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.planets ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.films ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.constants ALTER COLUMN id DROP DEFAULT;
ALTER TABLE public.characters ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE public.vehicles_id_seq;
DROP TABLE public.vehicles;
DROP SEQUENCE public.species_id_seq;
DROP TABLE public.species;
DROP SEQUENCE public.planets_id_seq;
DROP TABLE public.planets;
DROP SEQUENCE public.films_id_seq;
DROP TABLE public.films;
DROP SEQUENCE public.constants_id_seq;
DROP TABLE public.constants;
DROP SEQUENCE public.characters_id_seq;
DROP TABLE public.characters;
DROP FUNCTION public.notify_updated();
DROP FUNCTION public.notify_deleted();
DROP FUNCTION public.create_verify_normalized_relations();
DROP FUNCTION public.create_update_updated();
DROP FUNCTION public.create_insert_created();
--
-- Name: create_insert_created(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.create_insert_created() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
    declare
      _now timestamp := current_timestamp;
    begin
      new.updated = _now;
      new.created = _now;
      return new;
    end;
    $$;


--
-- Name: create_update_updated(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.create_update_updated() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
    begin
      new.updated = current_timestamp;
      return new;
    end;
    $$;


--
-- Name: create_verify_normalized_relations(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.create_verify_normalized_relations() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
    $$;


--
-- Name: notify_deleted(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.notify_deleted() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
    declare
      channel text := TG_ARGV[0];
    begin
      perform
        (
          with
            deleted     as (
              select
                id
              from old_table
              order by id
            ),
            stringified as (select json_agg(id)::text as deleted from deleted)
          select
            pg_notify(channel, (select deleted from stringified))
        );
      return null;
    end;
    $$;


--
-- Name: notify_updated(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.notify_updated() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
    declare
      channel text := TG_ARGV[0];
    begin
      perform
        (
          with
            updated     as (
              select
                id
              from new_table
              order by id
            ),
            stringified as (select json_agg(id)::text as updated from updated)
          select
            pg_notify(channel, (select updated from stringified))
        );
      return null;
    end;
    $$;


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: characters; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.characters (
    id bigint NOT NULL,
    name text NOT NULL,
    height integer,
    mass double precision,
    hair_color text[] DEFAULT ARRAY[]::text[],
    skin_color text[] DEFAULT ARRAY[]::text[],
    eye_color text,
    birth_year character varying(10),
    gender character varying(10),
    home_world bigint,
    species bigint,
    vehicles integer[] DEFAULT ARRAY[]::integer[],
    description text,
    created timestamp without time zone NOT NULL,
    updated timestamp without time zone NOT NULL
);


--
-- Name: characters_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.characters_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: characters_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.characters_id_seq OWNED BY public.characters.id;


--
-- Name: constants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.constants (
    id bigint NOT NULL,
    constant_type text NOT NULL,
    constant_value text NOT NULL,
    updated timestamp without time zone NOT NULL,
    created timestamp without time zone NOT NULL
);


--
-- Name: constants_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.constants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: constants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.constants_id_seq OWNED BY public.constants.id;


--
-- Name: films; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.films (
    id bigint NOT NULL,
    title text NOT NULL,
    episode integer NOT NULL,
    opening_crawl text NOT NULL,
    director text NOT NULL,
    producer text NOT NULL,
    release_date timestamp without time zone NOT NULL,
    characters integer[] DEFAULT ARRAY[]::integer[],
    star_ships integer[] DEFAULT ARRAY[]::integer[],
    vehicles integer[] DEFAULT ARRAY[]::integer[],
    species integer[] DEFAULT ARRAY[]::integer[],
    planets integer[] DEFAULT ARRAY[]::integer[],
    description text,
    created timestamp without time zone NOT NULL,
    updated timestamp without time zone NOT NULL
);


--
-- Name: films_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.films_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: films_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.films_id_seq OWNED BY public.films.id;


--
-- Name: planets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.planets (
    id bigint NOT NULL,
    name text NOT NULL,
    rotation_period double precision,
    orbital_period double precision,
    diameter double precision,
    climate text[] DEFAULT ARRAY[]::text[],
    terrain text[] DEFAULT ARRAY[]::text[],
    gravity double precision,
    surface_water double precision,
    population bigint,
    description text,
    created timestamp without time zone,
    updated timestamp without time zone
);


--
-- Name: planets_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.planets_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: planets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.planets_id_seq OWNED BY public.planets.id;


--
-- Name: species; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.species (
    id bigint NOT NULL,
    name text NOT NULL,
    classification text,
    average_height double precision,
    skin_colors text[] DEFAULT ARRAY[]::text[],
    hair_colors text[] DEFAULT ARRAY[]::text[],
    eye_colors text[] DEFAULT ARRAY[]::text[],
    average_lifespan double precision,
    home_world bigint,
    spoken_language text,
    description text,
    created timestamp without time zone NOT NULL,
    updated timestamp without time zone NOT NULL
);


--
-- Name: species_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.species_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: species_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.species_id_seq OWNED BY public.species.id;


--
-- Name: vehicles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vehicles (
    id bigint NOT NULL,
    name text NOT NULL,
    model text NOT NULL,
    manufacturer text,
    cost_in_credits bigint,
    size double precision,
    max_atmospheric_speed double precision,
    crew integer,
    passengers integer,
    cargo_capacity bigint,
    consumables interval,
    hyperdrive_rating integer,
    mglt integer,
    starship_class text,
    created timestamp without time zone NOT NULL,
    updated timestamp without time zone NOT NULL
);


--
-- Name: vehicles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.vehicles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vehicles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.vehicles_id_seq OWNED BY public.vehicles.id;


--
-- Name: characters id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.characters ALTER COLUMN id SET DEFAULT nextval('public.characters_id_seq'::regclass);


--
-- Name: constants id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.constants ALTER COLUMN id SET DEFAULT nextval('public.constants_id_seq'::regclass);


--
-- Name: films id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.films ALTER COLUMN id SET DEFAULT nextval('public.films_id_seq'::regclass);


--
-- Name: planets id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.planets ALTER COLUMN id SET DEFAULT nextval('public.planets_id_seq'::regclass);


--
-- Name: species id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.species ALTER COLUMN id SET DEFAULT nextval('public.species_id_seq'::regclass);


--
-- Name: vehicles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles ALTER COLUMN id SET DEFAULT nextval('public.vehicles_id_seq'::regclass);


--
-- Data for Name: characters; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.characters (id, name, height, mass, hair_color, skin_color, eye_color, birth_year, gender, home_world, species, vehicles, description, created, updated) FROM stdin;
1	Luke Skywalker	172	77	{BLONDE}	{FAIR}	BLUE	19BBY	male	1	1	{14,30,12,22}	Luke Skywalker is a fictional character and the main protagonist of the original film trilogy of the Star Wars franchise created by George Lucas. The character, portrayed by Mark Hamill, is an important figure in the Rebel Alliance's struggle against the Galactic Empire. He is the twin brother of Rebellion leader Princess Leia Organa of Alderaan, a friend and brother-in-law of smuggler Han Solo, an apprentice to Jedi Masters Obi-Wan "Ben" Kenobi and Yoda, the son of fallen Jedi Anakin Skywalker (Darth Vader) and Queen of Naboo/Republic Senator Padmé Amidala and maternal uncle of Kylo Ren / Ben Solo. The now non-canon Star Wars expanded universe depicts him as a powerful Jedi Master, husband of Mara Jade, the father of Ben Skywalker and maternal uncle of Jaina, Jacen and Anakin Solo.\nIn 2015, the character was selected by Empire magazine as the 50th greatest movie character of all time.[2] On their list of the 100 Greatest Fictional Characters, Fandomania.com ranked the character at number 14.[3]	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
2	C-3PO	167	75	\N	{GOLD}	YELLOW	112BBY	\N	1	2	{}	C-3PO (/siːˈθriːpi.oʊ/) or See-Threepio is a humanoid robot character from the Star Wars franchise who appears in the original trilogy, the prequel trilogy and the sequel trilogy. Built by Anakin Skywalker, C-3PO was designed as a protocol droid intended to assist in etiquette, customs, and translation, boasting that he is "fluent in over six million forms of communication". Along with his astromech droid counterpart and friend R2-D2, C-3PO provides comic relief within the narrative structure of the films, and serves as a foil. Anthony Daniels has portrayed the character in all nine Star Wars cinematic films released to date, including Rogue One and the animated The Clone Wars; C-3PO and R2-D2 are the only characters to appear in every film.\nDespite his oblivious nature, C-3PO has played a pivotal role in the Galaxy's history, appearing under the service of Shmi Skywalker, the Lars homestead, Padmé Amidala, Raymus Antilles, Luke Skywalker, and Leia Organa. In the majority of depictions, C-3PO's physical appearance is primarily a polished gold plating, although his appearance varies throughout the films; including the absence of metal coverings in The Phantom Menace, a dull copper plating in Attack of the Clones, a silver lower right leg introduced in A New Hope, and a red left arm in The Force Awakens.[1] C-3PO also appears frequently in the Star Wars Canon and Star Wars Legends continuities of novels, comic books, and video games, and was the protagonist in the ABC television show Droids.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
3	R2-D2	96	32	\N	{WHITE,BLUE}	RED	33BBY	\N	8	2	{}	R2-D2 is a fictional robot character in the Star Wars franchise created by George Lucas, who appears in the original trilogy, the prequel trilogy, the sequel trilogy, and Rogue One. A small astromech droid, R2-D2 is a major character and appears in all Star Wars films to date. Throughout the course of the films, he joins or supports Padmé Amidala, Anakin Skywalker, Leia Organa, Luke Skywalker, and Obi-Wan Kenobi in various points in the saga.\nEnglish actor Kenny Baker played R2-D2 in all three original Star Wars films, and received billing credit for the character in the prequel trilogy, where Baker's role was reduced, as R2-D2 was portrayed mainly by radio controlled props and CGI models. In the sequel trilogy, Baker was credited as consultant for The Force Awakens, however Jimmy Vee also co-performed the character in some scenes. Vee later took over the role beginning in The Last Jedi.[1] R2-D2's sounds and vocal effects were created by Ben Burtt. R2-D2 was designed in artwork by Ralph McQuarrie, co-developed by John Stears and built by Tony Dyson.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
4	Darth Vader	202	136	{NONE}	{WHITE}	YELLOW	41.9BBY	male	1	1	{13}	Darth Vader, also known by his birth name Anakin Skywalker, is a fictional character in the Star Wars franchise.[1][2][3] Vader appears in the original film trilogy as a pivotal antagonist whose actions drive the plot, while his past as Anakin Skywalker and the story of his corruption are central to the narrative of the prequel trilogy.\nThe character was created by George Lucas and has been portrayed by numerous actors. His appearances span the first six Star Wars films, as well as Rogue One, and his character is heavily referenced in Star Wars: The Force Awakens. He is also an important character in the Star Wars expanded universe of television series, video games, novels, literature and comic books. Originally a Jedi prophesied to bring balance to the Force, he falls to the dark side of the Force and serves the evil Galactic Empire at the right hand of his Sith master, Emperor Palpatine (also known as Darth Sidious).[4] He is also the father of Luke Skywalker and Princess Leia Organa, secret husband of Padmé Amidala and grandfather of Kylo Ren.\nDarth Vader has become one of the most iconic villains in popular culture, and has been listed among the greatest villains and fictional characters ever.[5][6] The American Film Institute listed him as the third greatest movie villain in cinema history on 100 Years... 100 Heroes and Villains, behind Hannibal Lecter and Norman Bates.[7] However, other critics consider him a tragic hero, citing his original motivations for the greater good before his fall to the dark side.[8][9]	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
5	Leia Organa	150	49	{BROWN}	{LIGHT}	BROWN	19BBY	female	2	1	{30}	Princess Leia Organa of Alderaan (also Senator Leia Organa or General Leia Organa) is a fictional character in the Star Wars franchise, portrayed in films by Carrie Fisher. Introduced in the original Star Wars film in 1977, Leia is princess of the planet Alderaan, a member of the Imperial Senate and an agent of the Rebel Alliance. She thwarts the sinister Sith Lord Darth Vader and helps bring about the destruction of the Empire's cataclysmic superweapon, the Death Star. In The Empire Strikes Back (1980), Leia commands a Rebel base and evades Vader as she falls in love with the smuggler, Han Solo. In Return of the Jedi (1983), Leia leads the operation to rescue Han from the crime lord Jabba the Hutt, and is revealed to be Vader's daughter and the twin sister of Luke Skywalker. The prequel film Revenge of the Sith (2005) establishes that the twins' mother is Senator (and former queen) Padmé Amidala of Naboo, who dies after childbirth. Leia is adopted by Senator Bail and Queen Breha Organa of Alderaan. In The Force Awakens (2015), Leia is the founder and General of the Resistance against the First Order and has a son with Han named Ben, who goes by the name Kylo Ren.\nIn the original Star Wars expanded universe (1977–2014) of novels, comics and video games, which are set in an alternate continuity, Leia continues her adventures with Han and Luke after Return of the Jedi, fighting Imperial resurgences and new threats to the galaxy. She becomes the Chief of State of the New Republic and a Jedi Master, and is the mother to three children by Han: Jaina, Jacen and Anakin Solo.\nOne of the more popular Star Wars characters, Leia has been called a 1980s icon, a feminist hero and model for other adventure heroines. She has appeared in many derivative works and merchandising, and has been referenced or parodied in several TV shows and films. Her "cinnamon buns" hairstyle from Star Wars (1977) and metal bikini from Return of the Jedi have become cultural icons.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
6	Owen Lars	178	120	{BROWN,GRAY}	{LIGHT}	BLUE	52BBY	male	1	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
7	Beru Whitesun lars	165	75	{BROWN}	{LIGHT}	BLUE	47BBY	female	1	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
8	R5-D4	97	32	\N	{WHITE,RED}	RED	\N	\N	1	2	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
9	Biggs Darklighter	183	84	{BLACK}	{LIGHT}	BROWN	24BBY	male	1	1	{12}	Rogue Squadron is a starfighter squadron in the Star Wars franchise. Many surviving members of Red Squadron, the X-wing attack force that Luke Skywalker joins during the Battle of Yavin in Star Wars Episode IV: A New Hope (1977), later join Rogue Squadron. The squadron appears in The Empire Strikes Back (1980) as Rogue Group. In the 2016 film Rogue One, Rebel fighters on a suicide mission to steal the plans for the Death Star self-identify as "Rogue One", a possible precursor to Rogue Squadron.[1]\nRogue Squadron is prominently featured in the comic book series Star Wars: X-wing, the ten-volume novel series Star Wars: X-wing, and the video game series Star Wars: Rogue Squadron. The unit is depicted as consisting of "the best pilots and the best fighters".[2]	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
10	Obi-Wan Kenobi	182	77	{AUBURN,WHITE}	{FAIR}	BLUE-GRAY	57BBY	male	20	1	{38,48,59,64,65,74}	Obi-Wan "Ben" Kenobi is a fictional character in the Star Wars franchise. Within the original trilogy he is portrayed by Sir Alec Guinness, while in the prequel trilogy a younger version of the character is portrayed by Ewan McGregor. In the original trilogy, he is a mentor to Luke Skywalker, to whom he introduces the ways of the Jedi. In the prequel trilogy, he is a master and friend to Anakin Skywalker. He is frequently featured as a main character in various other Star Wars media.\nSir Alec Guinness's portrayal of Obi-Wan in the original Star Wars (1977) remains the only time an actor has received an Oscar nomination (Best Supporting Actor) for acting in a Star Wars film.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
11	Anakin Skywalker	188	84	{BLOND}	{FAIR}	BLUE	41.9BBY	male	1	1	{44,46,59,65,39}	Darth Vader, also known by his birth name Anakin Skywalker, is a fictional character in the Star Wars franchise.[1][2][3] Vader appears in the original film trilogy as a pivotal antagonist whose actions drive the plot, while his past as Anakin Skywalker and the story of his corruption are central to the narrative of the prequel trilogy.\nThe character was created by George Lucas and has been portrayed by numerous actors. His appearances span the first six Star Wars films, as well as Rogue One, and his character is heavily referenced in Star Wars: The Force Awakens. He is also an important character in the Star Wars expanded universe of television series, video games, novels, literature and comic books. Originally a Jedi prophesied to bring balance to the Force, he falls to the dark side of the Force and serves the evil Galactic Empire at the right hand of his Sith master, Emperor Palpatine (also known as Darth Sidious).[4] He is also the father of Luke Skywalker and Princess Leia Organa, secret husband of Padmé Amidala and grandfather of Kylo Ren.\nDarth Vader has become one of the most iconic villains in popular culture, and has been listed among the greatest villains and fictional characters ever.[5][6] The American Film Institute listed him as the third greatest movie villain in cinema history on 100 Years... 100 Heroes and Villains, behind Hannibal Lecter and Norman Bates.[7] However, other critics consider him a tragic hero, citing his original motivations for the greater good before his fall to the dark side.[8][9]	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
12	Wilhuff Tarkin	180	\N	{AUBURN,GRAY}	{FAIR}	BLUE	64BBY	male	21	1	{}	Governor Wilhuff "Grand Moff" Tarkin, is a fictional character in the Star Wars franchise, first portrayed by Peter Cushing in the 1977 film Star Wars. He is the commander of the Death Star, the Galactic Empire's dwarf planet-sized super weapon. The character has been called "one of the most formidable villains in Star Wars history."[1]	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
13	Chewbacca	228	112	{BROWN}	\N	BLUE	200BBY	male	14	3	{19,10,22}	Chewbacca (/tʃuːˈbɑːkə/), nicknamed "Chewie", is a fictional character in the Star Wars franchise. He is a Wookiee, a tall, hirsute biped and intelligent species from the planet Kashyyyk. Chewbacca is the loyal friend and first mate of Han Solo, and serves as co-pilot on Solo's spaceship, the Millennium Falcon.[1] Within the films of the main saga, Chewbacca is portrayed by Peter Mayhew in the Star Wars on Episodes from III to VIII (Mayhew shares the role with his body double Joonas Suotamo on Episode VII and VIII). Suotamo took over the role alone in Solo: A Star Wars Story. The character has also appeared on television, books, comics, and video games.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
14	Han Solo	180	80	{BROWN}	{FAIR}	BROWN	where cons	male	22	1	{10,22}	Han Solo is a character in the Star Wars franchise. In the original film trilogy, Han and his co-pilot, Chewbacca, became involved in the Rebel Alliance's struggle against the Galactic Empire. During the course of the Star Wars narrative, he becomes a chief figure in the Alliance and succeeding galactic governments. Star Wars creator George Lucas described the character as "a loner who realizes the importance of being part of a group and helping for the common good".[2] Harrison Ford portrayed the character in the original Star Wars trilogy as well as The Force Awakens. Alden Ehrenreich will portray a young Han Solo in Solo: A Star Wars Story.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
15	Greedo	173	74	\N	{GREEN}	BLACK	44BBY	male	23	4	{}	Greedo (or Greedo the Young) is a fictional character from the Star Wars franchise. He is portrayed by Paul Blake as well as Maria De Aragon in some close-in pickup shots in Star Wars (1977). The character is part of a fan controversy known as "Han shot first".	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
31	Nien Nunb	160	68	{NONE}	{GRAY}	BLACK	\N	male	33	10	{10}	Nien Nunb is a fictional character in the Star Wars franchise. Introduced in the 1983 film Return of the Jedi, he was brought to life both as a puppet and a costumed actor during the film. Nunb was puppeteered by Mike Quinn and was portrayed by Richard Bonehill in wide shots. The character was voiced by Kipsang Rotich, a Kenyan student who spoke in his native Kalenjin language, as well as in the Kikuyu language. Quinn and Rotich returned to the role for the 2015 film Star Wars: The Force Awakens, with the former confirmed for the 2017 film Star Wars: The Last Jedi.\nWithin the fictional Star Wars universe, Nien Nunb was an arms dealer of the Sullustan species who joined the Alliance to Restore the Republic during the Galactic Civil War. Three decades later, Nunb was a starfighter pilot in the Resistance. He fought in the First Order–Resistance conflict, including the Battle of Starkiller Base.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
32	Qui-Gon Jinn	193	89	{BROWN}	{FAIR}	BLUE	92BBY	male	\N	1	{38}	Qui-Gon Jinn is a fictional character in the Star Wars franchise, portrayed by Liam Neeson as the main protagonist of the 1999 film Star Wars: Episode I – The Phantom Menace.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
16	Jabba Desilijic Tiure	175	1	\N	{GREEN,BROWN}	ORANGE	600BBY	hermaphrod	24	5	{}	Jabba Desilijic Tiure,[1] commonly known as Jabba the Hutt, is a fictional character and an antagonist in the Star Wars franchise created by George Lucas. He is depicted as a large, slug-like[2] alien. His appearance has been described by film critic Roger Ebert as a cross between a toad and the Cheshire Cat.[3]\nIn the original theatrical releases of the original Star Wars trilogy, Jabba the Hutt first appeared in Return of the Jedi (1983), though he is mentioned in Star Wars (1977) and The Empire Strikes Back (1980), and a previously deleted scene involving Jabba the Hutt was added to the 1997 theatrical re-release and subsequent home media releases of Star Wars. Jabba is introduced as the de facto leader of the Desilijic-Hutt Cartel, and the most powerful crime boss on Tatooine, who has a bounty on Han Solo's head. Jabba employs a retinue of career criminals, bounty hunters, smugglers, assassins and bodyguards to operate his criminal empire. He keeps a host of entertainers at his disposal at his palace: slaves, droids and alien creatures. Jabba has a grim sense of humor, an insatiable appetite, and affinities for gambling, slave girls, and torture.[1]\nThe character was incorporated into the Star Wars merchandising campaign that corresponded with the theatrical release of Return of the Jedi. Besides the films, Jabba the Hutt is featured in Star Wars Legends literature. Jabba the Hutt's image has since played an influential role in popular culture, particularly in the United States. The name is used as a satirical literary device and a political caricature to underscore negative qualities such as morbid obesity and corruption.[4][5]	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
18	Wedge Antilles	170	77	{BROWN}	{FAIR}	HAZEL	21BBY	male	22	1	{14,12}	Wedge Antilles is a fictional character in the Star Wars franchise. He is a supporting character portrayed by Denis Lawson in the original Star Wars trilogy.[1] Antilles is a starfighter pilot for the Rebel Alliance, and founded Rogue Squadron with his friend Luke Skywalker. Wedge is notable for being the only Rebel pilot to have survived both attack runs on the Death Stars at the battles of Yavin and Endor.[1] He also appears in the Star Wars expanded universe, most notably as the lead character in most of the X-Wing novels.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
19	Jek Tono Porkins	180	110	{BROWN}	{FAIR}	BLUE	\N	male	26	1	{12}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
20	Yoda	66	17	{WHITE}	{GREEN}	BROWN	896BBY	male	\N	6	{}	Yoda is a fictional character in the Star Wars franchise created by George Lucas, first appearing in the 1980 film The Empire Strikes Back. In the original films, he trains Luke Skywalker to fight against the Galactic Empire. In the prequel films, he serves as the Grand Master of the Jedi Order and as a high-ranking general of Clone Troopers in the Clone Wars. Following his death in Return of the Jedi at the age of 900, Yoda was the oldest living character in the Star Wars franchise in canon, until the introduction of Maz Kanata in Star Wars: The Force Awakens.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
21	Palpatine	170	75	{GRAY}	{PALE}	YELLOW	82BBY	male	8	1	{}	Sheev Palpatine[3], also known as Darth Sidious and The Emperor, is a fictional character and one of the primary antagonists of the Star Wars franchise,[4] mainly portrayed by Ian McDiarmid. In the original trilogy, he is depicted as the aged, pale-faced and cloaked Emperor of the Galactic Empire and the master of Darth Vader. In the prequel trilogy, he is portrayed as a charismatic Senator from Naboo who uses deception and political manipulation to rise to the position of Supreme Chancellor of the Galactic Republic, and then reorganizes the Republic into the Galactic Empire, with himself as Emperor.\nThough outwardly appearing to be a well-intentioned public servant and supporter of democracy prior to becoming Emperor,[5] he is actually Darth Sidious, the Dark Lord of the Sith – a cult of practitioners of the dark side of the Force previously thought to have been extinct in the Star Wars galaxy for a millennium.[5] As Sidious, he instigates the Clone Wars, nearly destroys the Jedi, and transforms the Republic into the Empire. He also manipulates Jedi Knight Anakin Skywalker into turning to the dark side and serving at his side as Darth Vader. Palpatine's reign is brought to an end when Vader kills him to save his son, Luke Skywalker.\nSince the initial theatrical run of Return of the Jedi, Palpatine has become a widely recognized popular culture symbol of evil, sinister deception, dictatorship, tyranny, and the subversion of democracy.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
22	Boba Fett	183	78	{BLACK}	{FAIR}	BROWN	31.5BBY	male	10	1	{21}	Boba Fett is a fictional character in the Star Wars franchise. In The Empire Strikes Back and Return of the Jedi, he is a bounty hunter hired by Darth Vader and also employed by Jabba the Hutt. He was also added briefly to the original film Star Wars when the film was digitally remastered. Star Wars: Episode II – Attack of the Clones establishes his origin as an unaltered clone of the bounty hunter Jango Fett raised as his son. He also appears in several episodes of Star Wars: The Clone Wars cartoon series which further describes his growth as a villain in the Star Wars universe. His aura of danger and mystery has created a cult following for the character.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
23	IG-88	200	140	{NONE}	{METAL}	RED	15BBY	none	\N	2	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
24	Bossk	190	113	{NONE}	{GREEN}	RED	53BBY	male	29	7	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
25	Lando Calrissian	177	79	{BLACK}	{DARK}	BROWN	31BBY	male	30	1	{10}	Lando Calrissian is a fictional character in the Star Wars franchise. He is portrayed by Billy Dee Williams in The Empire Strikes Back and Return of the Jedi, and will be played by Donald Glover in the upcoming standalone film, Solo: A Star Wars Story.[1] He also appears frequently in the Star Wars expanded universe of novels, comic books and video games, including a series of novels in which he is the protagonist.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
26	Lobot	175	79	{NONE}	{LIGHT}	BLUE	37BBY	male	6	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
27	Ackbar	180	83	{NONE}	{BROWN}	ORANGE	41BBY	male	31	8	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
28	Mon Mothma	150	\N	{AUBURN}	{FAIR}	BLUE	48BBY	female	32	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
29	Arvel Crynyd	\N	\N	{BROWN}	{FAIR}	BROWN	\N	male	\N	1	{28}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
30	Wicket Systri Warrick	88	20	{BROWN}	{BROWN}	BROWN	8BBY	male	7	9	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
33	Nute Gunray	191	90	{NONE}	{GREEN}	RED	\N	male	18	11	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
34	Finis Valorum	170	\N	{BLOND}	{FAIR}	BLUE	91BBY	male	9	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
35	Padmé Amidala	165	45	{BROWN}	{LIGHT}	BROWN	46BBY	female	8	1	{49,64,39}	Padmé Amidala (born Padmé Naberrie) is a fictional character in the Star Wars franchise, appearing in the prequel trilogy portrayed by actress Natalie Portman. She served as the Princess of Theed and later Queen of Naboo. After her reign, she became a senator in the Galactic Senate, an anti-war movement spokesperson, and co-founder of the opposition-faction that later emerged as the Rebel Alliance.[2] She was secretly married to the Jedi Anakin Skywalker, and the biological mother of Luke Skywalker and Leia Organa, which makes her the mother-in-law of Han Solo, and the grandmother of Kylo Ren.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
36	Jar Jar Binks	196	66	{NONE}	{ORANGE}	ORANGE	52BBY	male	8	12	{}	Jar Jar Binks is a fictional character from the Star Wars saga created by George Lucas. A major character in Star Wars: Episode I – The Phantom Menace, he also has a smaller role in Episode II: Attack of the Clones, and a one-line cameo in Episode III: Revenge of the Sith, as well as a role in the television series Star Wars: The Clone Wars. The first lead computer generated character of the franchise, he has been portrayed by Ahmed Best in most of his appearances.\nJar Jar's primary role in Episode I was to provide comic relief for the audience. Upon the movie's release, he was met with an overwhelmingly negative reception from both critics and audiences, and is today considered one of the most hated characters in the history of film.[1][2]	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
37	Roos Tarpals	224	82	{NONE}	{GRAY}	ORANGE	\N	male	8	12	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
38	Rugor Nass	206	\N	{NONE}	{GREEN}	ORANGE	\N	male	8	12	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
39	Ric Olié	183	\N	{BROWN}	{FAIR}	BLUE	\N	male	8	\N	{40}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
40	Watto	137	\N	{BLACK}	{BLUE,GRAY}	YELLOW	\N	male	34	13	{}	Watto is a fictional character in the Star Wars franchise, featured in the films The Phantom Menace and Attack of the Clones. He is computer-generated and played by voice actor Andy Secombe. He is a mean-tempered, greedy Toydarian, and owner of a second-hand goods store in Mos Espa on the planet Tatooine. Among Watto's belongings are the slaves Shmi Skywalker and her son, Anakin. He acquires them after winning a podracing bet with Gardulla the Hutt, and he puts them both to work in his store. Anakin demonstrates an incredible aptitude for equipment repair, and Watto decides to profit from it by having the boy fix various broken equipment in the store. He eventually loses Anakin in a podracing bet with Qui-Gon Jinn when he bets on a competitor, Sebulba, who is defeated by Anakin.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
41	Sebulba	112	40	{NONE}	{GRAY,RED}	ORANGE	\N	male	35	14	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
42	Quarsh Panaka	183	\N	{BLACK}	{DARK}	BROWN	62BBY	male	8	\N	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
43	Shmi Skywalker	163	\N	{BLACK}	{FAIR}	BROWN	72BBY	female	1	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
44	Darth Maul	175	80	{NONE}	{RED}	YELLOW	54BBY	male	36	22	{42,41}	Darth Maul, later known simply as Maul, is a fictional character in the Star Wars franchise. Trained as Darth Sidious's first apprentice, he serves as a Sith Lord and a master of wielding a double-bladed lightsaber. He first appears in Star Wars: Episode I – The Phantom Menace (portrayed by Ray Park and voiced by Peter Serafinowicz), and later makes appearances in Star Wars: The Clone Wars and Star Wars Rebels, voiced by Samuel Witwer.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
45	Bib Fortuna	180	\N	{NONE}	{PALE}	PINK	\N	male	37	15	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
46	Ayla Secura	178	55	{NONE}	{BLUE}	HAZEL	48BBY	female	37	15	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
47	Ratts Tyerell	79	15	{NONE}	{GRAY,BLUE}	\N	\N	male	38	16	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
48	Dud Bolt	94	45	{NONE}	{BLUE,GRAY}	YELLOW	\N	male	39	17	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
49	Gasgano	122	\N	{NONE}	{WHITE,BLUE}	BLACK	\N	male	40	18	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
50	Ben Quadinaros	163	65	{NONE}	{GRAY,GREEN,YELLOW}	ORANGE	\N	male	41	19	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
51	Mace Windu	188	84	{NONE}	{DARK}	BROWN	72BBY	male	42	1	{}	Mace Windu is a fictional character in the Star Wars franchise, portrayed by actor Samuel L. Jackson in the prequel films and voiced by voice-actor Terrence C. Carson in other projects. He appears as Master of the Jedi High Council and one of the last members of the order's upper echelons before the Galactic Republic's fall. He is the Council's primary liaison, although the Clone Wars caused him to question his most firmly held beliefs.[1]	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
52	Ki-Adi-Mundi	198	82	{WHITE}	{PALE}	YELLOW	92BBY	male	43	20	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
53	Kit Fisto	196	87	{NONE}	{GREEN}	BLACK	\N	male	44	21	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
54	Eeth Koth	171	\N	{BLACK}	{BROWN}	BROWN	\N	male	45	22	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
55	Adi Gallia	184	50	{NONE}	{DARK}	BLUE	\N	female	9	23	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
56	Saesee Tiin	188	\N	{NONE}	{PALE}	ORANGE	\N	male	47	24	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
57	Yarael Poof	264	\N	{NONE}	{WHITE}	YELLOW	\N	male	48	25	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
58	Plo Koon	188	80	{NONE}	{ORANGE}	BLACK	22BBY	male	49	26	{48}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
59	Mas Amedda	196	\N	{NONE}	{BLUE}	BLUE	\N	male	50	27	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
60	Gregar Typho	185	85	{BLACK}	{DARK}	BROWN	\N	male	8	1	{39}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
61	Cordé	157	\N	{BROWN}	{LIGHT}	BROWN	\N	female	8	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
62	Cliegg Lars	183	\N	{BROWN}	{FAIR}	BLUE	82BBY	male	1	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
63	Poggle the Lesser	183	80	{NONE}	{GREEN}	YELLOW	\N	male	11	28	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
64	Luminara Unduli	170	56	{BLACK}	{YELLOW}	BLUE	58BBY	female	51	29	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
65	Barriss Offee	166	50	{BLACK}	{YELLOW}	BLUE	40BBY	female	51	29	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
66	Dormé	165	\N	{BROWN}	{LIGHT}	BROWN	\N	female	8	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
67	Dooku	193	80	{WHITE}	{FAIR}	BROWN	102BBY	male	52	1	{55}	Count Dooku is a fictional character from the Star Wars franchise, appearing in Star Wars: Episode II – Attack of the Clones and Star Wars: Episode III – Revenge of the Sith (portrayed by Christopher Lee).[1] He was also voiced by Corey Burton in the animated series Star Wars: The Clone Wars and Star Wars: Clone Wars.\nOnce a respected Jedi Master, he falls to the dark side of the Force after the death of his former Padawan, Qui-Gon Jinn, and becomes Darth Sidious' second apprentice under the name Darth Tyranus. As the founder of the Confederacy of Independent Systems, he is instrumental in the Clone Wars. Dooku was trained by Yoda as a Padawan learner.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
68	Bail Prestor Organa	191	\N	{BLACK}	{TAN}	BROWN	67BBY	male	2	1	{}	Senator Bail Prestor Organa of Alderaan, is a fictional character in the Star Wars franchise, portrayed by actor Jimmy Smits in Attack of the Clones (2002), Revenge of the Sith (2005), and Rogue One (2016). He is the senator from the planet Alderaan, one of the founding members of the Rebel Alliance, and the adoptive father of Leia Organa, a main character in the franchise.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
69	Jango Fett	183	79	{BLACK}	{TAN}	BROWN	66BBY	male	53	1	{}	Jango Fett is a fictional character in the Star Wars franchise, created by George Lucas. He made his debut in the 2002 film Star Wars: Episode II – Attack of the Clones, where he was portrayed by actor Temuera Morrison.\nIn the context of the Star Wars universe, Jango Fett was regarded as the best mercenary in the galaxy during the final years of the Republic. A naturally skilled warrior, he was eventually chosen to serve as the genetic template for the Clone Army of the Galactic Republic. These clone soldiers were genetically modified to be predisposed toward unquestioning obedience to the chain of command, unlike their highly independent progenitor. Jango was also the "father" of unaltered clone Boba Fett, which he requested as part of his contract with the Kaminoans.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
70	Zam Wesell	168	55	{BLONDE}	{FAIR,GREEN,YELLOW}	YELLOW	\N	female	54	30	{45}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
71	Dexter Jettster	198	102	{NONE}	{BROWN}	YELLOW	\N	male	55	31	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
72	Lama Su	229	88	{NONE}	{GRAY}	BLACK	\N	male	10	32	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
73	Taun We	213	\N	{NONE}	{GRAY}	BLACK	\N	female	10	32	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
74	Jocasta Nu	167	\N	{WHITE}	{FAIR}	BLUE	\N	female	9	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
75	R4-P17	96	\N	{NONE}	{SILVER,RED}	RED, BLUE	\N	female	\N	\N	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
76	Wat Tambor	193	48	{NONE}	{GREEN,GRAY}	\N	\N	male	56	33	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
77	San Hill	191	\N	{NONE}	{GRAY}	GOLD	\N	male	57	34	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
78	Shaak Ti	178	57	{NONE}	{RED,BLUE,WHITE}	BLACK	\N	female	58	35	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
79	Grievous	216	159	{NONE}	{BROWN,WHITE}	GREEN, YELLOW	\N	male	59	36	{60,74}	General Grievous is a fictional character in the Star Wars franchise. A former Kaleesh warlord named Qymaen jai Sheelal, he is the Supreme Commander of the Confederacy of Independent Systems during the Clone Wars against the Galactic Republic and is trained in all lightsaber combat forms to ensure the Jedi's destruction.\nThe character was originally introduced in 2004 in the animated series Star Wars: Clone Wars, before being part of the film Star Wars: Episode III – Revenge of the Sith.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
80	Tarfful	234	136	{BROWN}	{BROWN}	BLUE	\N	male	14	3	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
81	Raymus Antilles	188	79	{BROWN}	{LIGHT}	BROWN	\N	male	2	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
82	Sly Moore	178	48	{NONE}	{PALE}	WHITE	\N	female	60	\N	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
83	Tion Medon	206	80	{NONE}	{GRAY}	BLACK	\N	male	12	37	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
84	Finn	\N	\N	{BLACK}	{DARK}	DARK	\N	male	\N	1	{}	Finn, designation number FN-2187, is a fictional character in the Star Wars franchise. He first appeared in the 2015 film Star Wars: The Force Awakens in which he is a stormtrooper for the First Order who flees and turns against it after being shocked by their cruelty in his first combat mission. He is played by British actor John Boyega, who will reprise the role in Star Wars: The Last Jedi.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
85	Rey	\N	\N	{BROWN}	{LIGHT}	HAZEL	\N	female	\N	1	{}	\N	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
86	Poe Dameron	\N	\N	{BROWN}	{LIGHT}	BROWN	\N	male	\N	1	{77}	Poe Dameron is a fictional character in the Star Wars franchise. Introduced in the 2015 film Star Wars: The Force Awakens, he is portrayed by Oscar Isaac. Poe is an X-wing fighter pilot for the Resistance who inadvertently brings renegade stormtrooper Finn (John Boyega) and Jakku scavenger Rey (Daisy Ridley) into the fight against—and eventually a victory over—the sinister First Order. He is featured in The Force Awakens media and merchandising as well as an eponymous comic book series, and will appear in the film's forthcoming sequel, Star Wars: The Last Jedi. Isaac and the character have received positive reviews, with Poe being compared to the characterization of Han Solo (Harrison Ford) in the original Star Wars film trilogy.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
87	BB8	\N	\N	{NONE}	{NONE}	BLACK	\N	none	\N	2	{}	BB-8 (or Beebee-Ate) is a droid character in the Star Wars franchise, first appearing in the 2015 film Star Wars: The Force Awakens. Spherical with a free-moving domed head, BB-8 is portrayed by both a rod puppet and a remote-controlled robotic unit.	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
88	Captain Phasma	\N	\N	\N	\N	\N	\N	female	\N	\N	{}	Captain Phasma is a fictional character in the Star Wars franchise, portrayed by Gwendoline Christie. Introduced in Star Wars: The Force Awakens (2015), the first film in the Star Wars sequel trilogy, Phasma is the commander of the First Order's force of stormtroopers. Christie confirmed that the character would reappear in the next of the trilogy's films, Star Wars: The Last Jedi. The character also made an additional appearance in Before the Awakening, an anthology book set before the events of The Force Awakens.\nJ. J. Abrams created Phasma from an armor design originally developed for Kylo Ren and named her after the 1979 horror film Phantasm. The character was originally conceived as male before being changed to female. Phasma appeared prominently in promotion and marketing for The Force Awakens, but the character's ultimately minor role in the film was the subject of criticism. Nonetheless, merchandise featuring the character found success and her figure was the bestselling of all Force Awakens action figures on Amazon.co.uk.[3]	2019-03-21 19:14:47.000687	2019-03-21 19:14:47.000687
\.


--
-- Data for Name: constants; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.constants (id, constant_type, constant_value, updated, created) FROM stdin;
1	TERRAIN	DESERTS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
2	TERRAIN	ROCKY_DESERTS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
3	TERRAIN	GRASSLANDS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
4	TERRAIN	MOUNTAINS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
5	TERRAIN	JUNGLES	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
6	TERRAIN	RAINFORESTS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
7	TERRAIN	TUNDRA	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
8	TERRAIN	ICE_CAVES	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
9	TERRAIN	SWAMPS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
10	TERRAIN	GAS_GIANT	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
11	TERRAIN	FORESTS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
12	TERRAIN	LAKES	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
13	TERRAIN	GRASSY_HILLS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
14	TERRAIN	CITYSCAPE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
15	TERRAIN	OCEAN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
16	TERRAIN	ROCK	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
17	TERRAIN	SCRUBLANDS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
18	TERRAIN	SAVANNA	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
19	TERRAIN	CANYONS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
20	TERRAIN	SINKHOLES	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
21	TERRAIN	VOLCANOES	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
22	TERRAIN	LAVA_RIVERS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
23	TERRAIN	CAVES	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
24	TERRAIN	RIVERS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
25	TERRAIN	AIRLESS_ASTEROID	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
26	TERRAIN	GLACIERS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
27	TERRAIN	ICE_CANYONS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
28	TERRAIN	FUNGUS_FORESTS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
29	TERRAIN	ROCK_ARCHES	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
30	TERRAIN	PLAINS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
31	TERRAIN	URBAN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
32	TERRAIN	HILLS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
33	TERRAIN	OCEANS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
34	TERRAIN	BOGS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
35	TERRAIN	SAVANNAS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
36	TERRAIN	ISLANDS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
37	TERRAIN	SEAS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
38	TERRAIN	MESAS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
39	TERRAIN	REEFS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
40	TERRAIN	ROCKY	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
41	TERRAIN	VALLEYS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
42	TERRAIN	BARREN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
43	TERRAIN	ASH	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
44	TERRAIN	TOXIC_CLOUDSEA	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
45	TERRAIN	PLATEAUS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
46	TERRAIN	VERDANT	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
47	TERRAIN	ROCKY_CANYONS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
48	TERRAIN	ACID_POOLS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
49	TERRAIN	VINES	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
50	TERRAIN	CITIES	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
51	TERRAIN	SAVANNAHS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
52	TERRAIN	CLIFFS	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
53	CLIMATE	ARID	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
54	CLIMATE	TEMPERATE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
55	CLIMATE	TROPICAL	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
56	CLIMATE	FROZEN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
57	CLIMATE	MURKY	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
58	CLIMATE	WINDY	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
59	CLIMATE	HOT	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
60	CLIMATE	ARTIFICIAL	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
61	CLIMATE	FRIGID	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
62	CLIMATE	HUMID	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
63	CLIMATE	MOIST	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
64	CLIMATE	POLLUTED	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
65	CLIMATE	SUPERHEATED	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
66	CLIMATE	ARCTIC	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
67	CLIMATE	SUBARCTIC	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
68	SKIN_COLOR	GOLD	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
69	SKIN_COLOR	WHITE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
70	SKIN_COLOR	BLUE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
71	SKIN_COLOR	LIGHT	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
72	SKIN_COLOR	RED	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
73	SKIN_COLOR	GREEN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
74	SKIN_COLOR	BROWN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
75	SKIN_COLOR	PALE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
76	SKIN_COLOR	METAL	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
77	SKIN_COLOR	DARK	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
78	SKIN_COLOR	GRAY	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
79	SKIN_COLOR	ORANGE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
80	SKIN_COLOR	YELLOW	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
81	SKIN_COLOR	TAN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
82	SKIN_COLOR	SILVER	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
83	SKIN_COLOR	NONE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
84	SKIN_COLOR	CAUCASIAN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
85	SKIN_COLOR	BLACK	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
86	SKIN_COLOR	ASIAN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
87	SKIN_COLOR	HISPANIC	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
88	SKIN_COLOR	MAGENTA	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
89	SKIN_COLOR	PURPLE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
90	SKIN_COLOR	PINK	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
91	SKIN_COLOR	PALE_PINK	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
92	SKIN_COLOR	PEACH	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
93	HAIR_COLOR	NONE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
94	HAIR_COLOR	BROWN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
95	HAIR_COLOR	BROWN, GRAY	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
96	HAIR_COLOR	BLACK	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
97	HAIR_COLOR	AUBURN, WHITE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
98	HAIR_COLOR	BLOND	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
99	HAIR_COLOR	AUBURN, GRAY	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
100	HAIR_COLOR	WHITE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
101	HAIR_COLOR	GRAY	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
102	HAIR_COLOR	AUBURN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
103	HAIR_COLOR	BLONDE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
104	HAIR_COLOR	RED	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
105	EYE_COLOR	BROWN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
106	EYE_COLOR	BLUE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
107	EYE_COLOR	GREEN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
108	EYE_COLOR	HAZEL	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
109	EYE_COLOR	GRAY	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
110	EYE_COLOR	AMBER	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
111	EYE_COLOR	YELLOW	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
112	EYE_COLOR	GOLDEN	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
113	EYE_COLOR	RED	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
114	EYE_COLOR	BLACK	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
115	EYE_COLOR	ORANGE	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
116	EYE_COLOR	PINK	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
117	EYE_COLOR	INDIGO	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
118	EYE_COLOR	SILVER	2019-03-21 19:13:56.856893	2019-03-21 19:13:56.856893
\.


--
-- Data for Name: films; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.films (id, title, episode, opening_crawl, director, producer, release_date, characters, star_ships, vehicles, species, planets, description, created, updated) FROM stdin;
1	A New Hope	4	It is a period of civil war.\nRebel spaceships, striking\nfrom a hidden base, have won\ntheir first victory against\nthe evil Galactic Empire.\nDuring the battle, Rebel\nspies managed to steal secret\nplans to the Empire's\nultimate weapon, the DEATH\nSTAR, an armored space\nstation with enough power\nto destroy an entire planet.\nPursued by the Empire's\nsinister agents, Princess\nLeia races home aboard her\nstarship, custodian of the\nstolen plans that can save her\npeople and restore\nfreedom to the galaxy....	George Lucas	Gary Kurtz, Rick McCallum	1977-05-25 00:00:00	{1,2,3,4,5,6,7,8,9,10,12,13,14,15,16,18,19,81}	{2,3,5,9,10,11,12,13}	{4,6,7,8}	{5,3,2,1,4}	{2,3,1}	Star Wars (later retitled Star Wars: Episode IV – A New Hope) is a 1977 American epic space opera film written and directed by George Lucas. It is the first film in the original Star Wars trilogy, the first Star Wars movie in general, and the beginning of the Star Wars franchise. Starring Mark Hamill, Harrison Ford, Carrie Fisher, Peter Cushing, Alec Guinness, David Prowse, James Earl Jones, Anthony Daniels, Kenny Baker, and Peter Mayhew, the film's plot focuses on the Rebel Alliance, led by Princess Leia (Fisher), and its attempt to destroy the Galactic Empire's space station, the Death Star. This conflict disrupts the isolated life of farmhand Luke Skywalker (Hamill), who inadvertently acquires a pair of droids that possess stolen architectural plans for the Death Star. When the Empire begins a destructive search for the missing droids, Skywalker accompanies Jedi Master Obi-Wan Kenobi (Guinness) on a mission to return the plans to the Rebel Alliance and rescue Leia from her imprisonment by the Empire.\nStar Wars was released in theatres in the United States on May 25, 1977. It earned $461 million in the U.S. and $314 million overseas, totaling $775 million ($3.132 billion in 2017 dollars, adjusted for inflation). It surpassed Jaws (1975) to become the highest-grossing film of all time until the release of E.T. the Extra-Terrestrial (1982). When adjusted for inflation, Star Wars is the second-highest-grossing film in North America, and the third-highest-grossing film in the world. It received ten Academy Award nominations (including Best Picture), winning seven. It was among the first films to be selected as part of the U.S. Library of Congress' National Film Registry as being "culturally, historically, or aesthetically significant". At the time, it was the most recent film on the registry and the only one chosen from the 1970s. Its soundtrack was added to the U.S. National Recording Registry in 2004. Today, it is often regarded as one of the most important films in the history of motion pictures. It launched an industry of tie-in products, including TV series spinoffs, novels, comic books, and video games, and merchandise including toys, games and clothing.\nThe film's success led to two critically and commercially successful sequels, The Empire Strikes Back in 1980 and Return of the Jedi in 1983. Star Wars was reissued multiple times at Lucas' behest, incorporating many changes including modified computer-generated effects, altered dialogue, re-edited shots, remixed soundtracks, and added scenes. A prequel trilogy was released beginning with The Phantom Menace in 1999, continuing with Attack of the Clones in 2002, and concluding with Revenge of the Sith in 2005. The film was followed by a sequel trilogy beginning with The Force Awakens in 2015. A direct prequel, Rogue One, was released in 2016. The film's fourth sequel, The Last Jedi, will be released in December 2017.	2019-03-21 19:15:00.6504	2019-03-21 19:15:00.6504
2	The Empire Strikes Back	5	It is a dark time for the\nRebellion. Although the Death\nStar has been destroyed,\nImperial troops have driven the\nRebel forces from their hidden\nbase and pursued them across\nthe galaxy.\nEvading the dreaded Imperial\nStarfleet, a group of freedom\nfighters led by Luke Skywalker\nhas established a new secret\nbase on the remote ice world\nof Hoth.\nThe evil lord Darth Vader,\nobsessed with finding young\nSkywalker, has dispatched\nthousands of remote probes into\nthe far reaches of space....	Irvin Kershner	Gary Kutz, Rick McCallum	1980-05-17 00:00:00	{1,2,3,4,5,10,13,14,18,20,21,22,23,24,25,26}	{10,11,12,15,21,22,23,3,17}	{8,14,16,18,19,20}	{6,7,3,2,1}	{4,5,6,27}	The Empire Strikes Back (also known as Star Wars: Episode V – The Empire Strikes Back) is a 1980 American epic space opera film directed by Irvin Kershner. Leigh Brackett and Lawrence Kasdan wrote the screenplay, with George Lucas writing the film's story and serving as executive producer. The second installment in the original Star Wars trilogy, it was produced by Gary Kurtz for Lucasfilm and stars Mark Hamill, Harrison Ford, Carrie Fisher, Billy Dee Williams, Anthony Daniels, David Prowse, Kenny Baker, Peter Mayhew, and Frank Oz.\nThe film is set three years after Star Wars. The Galactic Empire, under the leadership of the villainous Darth Vader and the Emperor, is in pursuit of Luke Skywalker and the rest of the Rebel Alliance. While Vader relentlessly pursues the small band of Luke's friends—Han Solo, Princess Leia Organa, and others—across the galaxy, Luke studies the Force under Jedi Master Yoda. When Vader captures Luke's friends, Luke must decide whether to complete his training and become a full Jedi Knight or to confront Vader and save them.\nFollowing a difficult production, The Empire Strikes Back was released on May 21, 1980. It received mixed reviews from critics initially but has since grown in esteem, becoming the most critically acclaimed chapter in the Star Wars saga; it is now widely regarded as one of the greatest films of all time.[5][6][7][8] The film ranks #3 on Empire's 2008 list of the 500 greatest movies of all time.[9] It became the highest-grossing film of 1980 and, to date, has earned more than $538 million worldwide from its original run and several re-releases. When adjusted for inflation, it is the second-highest-grossing sequel of all time and the 13th-highest-grossing film in North America.[10] The film was followed by Return of the Jedi, which was released in 1983.\nIn 2010, the film was selected for preservation in the United States' National Film Registry by the Library of Congress for being "culturally, historically, and aesthetically significant."	2019-03-21 19:15:00.6504	2019-03-21 19:15:00.6504
3	Return of the Jedi	6	Luke Skywalker has returned to\nhis home planet of Tatooine in\nan attempt to rescue his\nfriend Han Solo from the\nclutches of the vile gangster\nJabba the Hutt.\nLittle does Luke know that the\nGALACTIC EMPIRE has secretly\nbegun construction on a new\narmored space station even\nmore powerful than the first\ndreaded Death Star.\nWhen completed, this ultimate\nweapon will spell certain doom\nfor the small band of rebels\nstruggling to restore freedom\nto the galaxy...	Richard Marquand	Howard G. Kazanjian, George Lucas, Rick McCallum	1983-05-25 00:00:00	{1,2,3,4,5,10,13,14,16,18,20,21,22,25,27,28,29,30,31,45}	{10,11,12,15,22,23,27,28,29,3,17,2}	{8,16,18,19,24,25,26,30}	{5,6,8,9,10,15,3,2,1}	{5,7,8,9,1}	Return of the Jedi (also known as Star Wars: Episode VI – Return of the Jedi) is a 1983 American epic space opera film directed by Richard Marquand. The screenplay by Lawrence Kasdan and George Lucas was from a story by Lucas, who was also the executive producer. It is the third installment in the original Star Wars trilogy and the first film to use THX technology. The film is set one year after The Empire Strikes Back[9] and was produced by Howard Kazanjian for Lucasfilm Ltd. The film stars Mark Hamill, Harrison Ford, Carrie Fisher, Billy Dee Williams, Anthony Daniels, David Prowse, Kenny Baker, Peter Mayhew and Frank Oz.\nThe Galactic Empire, under the direction of the ruthless Emperor, is constructing a second Death Star in order to crush the Rebel Alliance once and for all. Since the Emperor plans to personally oversee the final stages of its construction, the Rebel Fleet launches a full-scale attack on the Death Star in order to prevent its completion and kill the Emperor, effectively bringing an end to the Empire's hold over the galaxy. Meanwhile, Luke Skywalker, a Jedi apprentice, struggles to bring his father Darth Vader back to the Light Side of the Force.\nDavid Lynch and David Cronenberg were considered to direct the project before Marquand signed on as director. The production team relied on Lucas' storyboards during pre-production. While writing the shooting script, Lucas, Kasdan, Marquand, and producer Howard Kazanjian spent two weeks in conference discussing ideas to construct it. Kazanjian's schedule pushed shooting to begin a few weeks early to allow Industrial Light & Magic more time to work on the film's effects in post-production. Filming took place in England, California, and Arizona from January to May 1982 (1982-05). Strict secrecy surrounded the production and the film used the working title Blue Harvest to prevent price gouging.\nThe film was released in theaters on May 25, 1983, six years to the day after the release of the first film, receiving mostly positive reviews. The film grossed between $475 million[6][7] and $572 million worldwide.[8] Several home video and theatrical releases and revisions to the film followed over the next 20 years. Star Wars continued with The Phantom Menace as part of the film series' prequel trilogy.\nA sequel, The Force Awakens, was released on December 18, 2015, as part of the new sequel trilogy.[10]	2019-03-21 19:15:00.6504	2019-03-21 19:15:00.6504
4	The Phantom Menace	1	Turmoil has engulfed the\nGalactic Republic. The taxation\nof trade routes to outlying star\nsystems is in dispute.\nHoping to resolve the matter\nwith a blockade of deadly\nbattleships, the greedy Trade\nFederation has stopped all\nshipping to the small planet\nof Naboo.\nWhile the Congress of the\nRepublic endlessly debates\nthis alarming chain of events,\nthe Supreme Chancellor has\nsecretly dispatched two Jedi\nKnights, the guardians of\npeace and justice in the\ngalaxy, to settle the conflict....	George Lucas	Rick McCallum	1999-05-19 00:00:00	{2,3,10,11,16,20,21,32,33,34,36,37,38,39,40,41,42,43,44,46,48,49,50,51,52,53,54,55,56,57,58,59,47,35}	{40,41,31,32,39}	{33,34,35,36,37,38,42}	{1,2,6,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27}	{8,9,1}	\N	2019-03-21 19:15:00.6504	2019-03-21 19:15:00.6504
5	Attack of the Clones	2	There is unrest in the Galactic\nSenate. Several thousand solar\nsystems have declared their\nintentions to leave the Republic.\nThis separatist movement,\nunder the leadership of the\nmysterious Count Dooku, has\nmade it difficult for the limited\nnumber of Jedi Knights to maintain \npeace and order in the galaxy.\nSenator Amidala, the former\nQueen of Naboo, is returning\nto the Galactic Senate to vote\non the critical issue of creating\nan ARMY OF THE REPUBLIC\nto assist the overwhelmed\nJedi....	George Lucas	Rick McCallum	2002-05-16 00:00:00	{2,3,6,7,10,11,20,21,22,33,36,40,43,46,51,52,53,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,82,35}	{21,39,43,47,48,49,32,52,58}	{4,44,45,46,50,51,53,54,55,56,57}	{32,33,2,35,6,1,12,34,13,15,28,29,30,31}	{8,9,10,11,1}	Star Wars: Episode II – Attack of the Clones is a 2002 American epic space opera film directed by George Lucas and written by Lucas and Jonathan Hales. It is the second installment of the Star Wars prequel trilogy, and stars Ewan McGregor, Natalie Portman, Hayden Christensen, Ian McDiarmid, Samuel L. Jackson, Christopher Lee, Temuera Morrison, Anthony Daniels, Kenny Baker and Frank Oz.\nSet ten years after the events in The Phantom Menace, the galaxy is on the brink of civil war. Led by a former Jedi named Count Dooku, thousands of planetary systems threaten to secede from the Galactic Republic. After Senator Padmé Amidala evades assassination, Jedi apprentice Anakin Skywalker becomes her protector, while his mentor Obi-Wan Kenobi investigates the attempt on Padmé's life. Soon Anakin, Padmé and Obi-Wan witness the onset of a new threat to the galaxy, the Clone Wars.\nDevelopment of Attack of the Clones began in March 2000, after the release of The Phantom Menace. By June 2000, Lucas and Hales completed a draft of the script and principal photography took place from June to September 2000. The film crew primarily shot at Fox Studios Australia in Sydney, Australia, with additional footage filmed in Tunisia, Spain and Italy. It was one of the first motion pictures shot completely on a high-definition digital 24-frame system.\nThe film was released in the United States on May 16, 2002. Some critics hailed it as an improvement over The Phantom Menace, while others called it the worst entry of the franchise.[5][6] The visual effects, costume design, musical score, action scenes and Ewan McGregor’s performance as Obi-Wan Kenobi were all praised, however the romance of Anakin and Padmé, the dialogue and the film's long runtime were all criticized. The film was a financial success, making over $649 million worldwide; however, it also became the first Star Wars film to be outgrossed in its year of release, placing third domestically and fourth internationally. The film was released on VHS and DVD on November 12, 2002 and was later released on Blu-ray on September 16, 2011. Following Attack of the Clones, the third and final film of the prequel trilogy, Revenge of the Sith, premiered in 2005.	2019-03-21 19:15:00.6504	2019-03-21 19:15:00.6504
6	Revenge of the Sith	3	War! The Republic is crumbling\nunder attacks by the ruthless\nSith Lord, Count Dooku.\nThere are heroes on both sides.\nEvil is everywhere.\nIn a stunning move, the\nfiendish droid leader, General\nGrievous, has swept into the\nRepublic capital and kidnapped\nChancellor Palpatine, leader of\nthe Galactic Senate.\nAs the Separatist Droid Army\nattempts to flee the besieged\ncapital with their valuable\nhostage, two Jedi Knights lead a\ndesperate mission to rescue the\ncaptive Chancellor....	George Lucas	Rick McCallum	2005-05-19 00:00:00	{1,2,3,4,5,6,7,10,11,12,13,20,21,33,46,51,52,53,54,55,56,58,63,64,67,68,75,78,79,80,81,82,83,35}	{48,59,61,32,63,64,65,66,68,74,75,2}	{33,50,60,62,67,69,70,71,72,73,76,53,56}	{19,33,2,3,36,37,6,1,34,15,35,20,23,24,25,26,27,28,29,30}	{2,5,8,9,12,13,14,15,16,17,18,19,1}	Star Wars: Episode III – Revenge of the Sith is a 2005 American epic space opera film written and directed by George Lucas. It is the sixth entry of the Star Wars film series and stars Ewan McGregor, Natalie Portman, Hayden Christensen, Ian McDiarmid, Samuel L. Jackson, Christopher Lee, Anthony Daniels, Kenny Baker, and Frank Oz. A sequel to The Phantom Menace (1999) and Attack of the Clones (2002), the film is the third and final installment of the Star Wars prequel trilogy.\nThe film begins three years after the onset of the Clone Wars. The Jedi Knights are spread across the galaxy, leading a massive war against the Separatists. The Jedi Council dispatches Jedi Master Obi-Wan Kenobi to eliminate the notorious General Grievous, leader of the Separatist Army. Meanwhile, Jedi Knight Anakin Skywalker grows close to Palpatine, the Supreme Chancellor of the Galactic Republic and, unknown to the public, a Sith Lord. Their deepening friendship threatens the Jedi Order, the Republic, and Anakin himself.\nLucas began writing the script before production of Attack of the Clones ended. Production of Revenge of the Sith started in September 2003, and filming took place in Australia with additional locations in Thailand, Switzerland, China, Italy and the United Kingdom. Revenge of the Sith premiered on May 15, 2005, at the Cannes Film Festival, then released worldwide on May 19, 2005. The film received generally positive reviews from critics, especially in contrast to the less positive reviews of the previous two prequels, including praise for its storyline, action scenes, John Williams' musical score, visual effects, and performances from Ewan McGregor, Ian McDiarmid, Frank Oz, and Jimmy Smits. It is the final film in the Star Wars franchise to be distributed by 20th Century Fox before The Walt Disney Company's acquisition of Lucasfilm in 2012.\nRevenge of the Sith broke several box office records during its opening week and went on to earn over $848 million worldwide,[4] making it, at the time, the third-highest-grossing film in the Star Wars franchise, unadjusted for inflation. It was the highest-grossing film of 2005 in the U.S. and the second-highest-grossing film of 2005 behind Harry Potter and the Goblet of Fire.[4] The film also holds the record for the highest opening day gross on a Thursday, making $50,013,859.00.[5] The Star Wars saga continued with the release of The Force Awakens, the first installment of the sequel trilogy, in 2015.[6][7]	2019-03-21 19:15:00.6504	2019-03-21 19:15:00.6504
7	The Force Awakens	7	Luke Skywalker has vanished.\nIn his absence, the sinister\nFIRST ORDER has risen from\nthe ashes of the Empire\nand will not rest until\nSkywalker, the last Jedi,\nhas been destroyed.\n \nWith the support of the\nREPUBLIC, General Leia Organa\nleads a brave RESISTANCE.\nShe is desperate to find her\nbrother Luke and gain his\nhelp in restoring peace and\njustice to the galaxy.\n \nLeia has sent her most daring\npilot on a secret mission\nto Jakku, where an old ally\nhas discovered a clue to\nLuke's whereabouts....	J. J. Abrams	Kathleen Kennedy, J. J. Abrams, Bryan Burk	2015-12-11 00:00:00	{1,3,5,13,14,27,84,85,86,87,88}	{77,10}	{}	{3,2,1}	{61}	Star Wars: The Force Awakens (also known as Star Wars: Episode VII – The Force Awakens) is a 2015 American epic space opera film co-written, co-produced and directed by J. J. Abrams. The sequel to 1983's Return of the Jedi, The Force Awakens is the first installment of the Star Wars sequel trilogy. It stars Harrison Ford, Mark Hamill, Carrie Fisher, Adam Driver, Daisy Ridley, John Boyega, Oscar Isaac, Lupita Nyong'o, Andy Serkis, Domhnall Gleeson, Anthony Daniels, Peter Mayhew, and Max von Sydow. Produced by Lucasfilm Ltd. and Abrams' production company Bad Robot Productions and distributed worldwide by Walt Disney Studios Motion Pictures, The Force Awakens was the first Star Wars film not produced by franchise creator George Lucas. Set 30 years after Return of the Jedi, it follows Rey, Finn and Poe Dameron's search for Luke Skywalker and their fight alongside the Resistance, led by veterans of the Rebel Alliance, against Kylo Ren and the First Order, a successor to the Galactic Empire.\nThe Force Awakens was announced after The Walt Disney Company's acquisition of Lucasfilm in October 2012. It was produced by Abrams, his longtime collaborator Bryan Burk, and Lucasfilm president Kathleen Kennedy. Abrams and Lawrence Kasdan, co-writer of the original trilogy films The Empire Strikes Back (1980) and Return of the Jedi (1983), rewrote an initial script by Michael Arndt. John Williams, composer for the previous six films, returned to compose the film's score. Lucas served as creative consultant during the film's early production. Filming began in April 2014 in Abu Dhabi and Iceland, with principal photography also taking place in Ireland and Pinewood Studios in the United Kingdom, and concluded in November 2014. It is the first live-action Star Wars film since Revenge of the Sith (2005).\nStar Wars: The Force Awakens was widely anticipated, and Disney backed the film with extensive marketing campaigns. It premiered in Los Angeles on December 14, 2015, four days before its wide release. The film received positive reviews, with its ensemble cast, direction, musical score, visual effects, and action sequences receiving particular praise, though it received some criticism for being derivative of the original trilogy. The film broke various box office records and became, unadjusted for inflation, the highest-grossing installment in the franchise, the highest-grossing film in North America, and the third-highest-grossing film of all time, with a worldwide gross of over $2 billion and a net profit of over $780 million.[6][7] It received five Academy Award nominations and four British Academy Film Award nominations, where it won Best Special Visual Effects. Two sequels, The Last Jedi and Episode IX, are scheduled for release in 2017 and 2019, respectively.	2019-03-21 19:15:00.6504	2019-03-21 19:15:00.6504
\.


--
-- Data for Name: planets; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.planets (id, name, rotation_period, orbital_period, diameter, climate, terrain, gravity, surface_water, population, description, created, updated) FROM stdin;
1	Tatooine	23	304	10465	{ARID}	{DESERTS}	1	1	200000	Tatooine /ˌtætuːˈiːn/ is a fictional desert planet that appears in the Star Wars space opera franchise. It is beige-coloured and is depicted as a remote, desolate world orbiting a pair of binary stars, and inhabited by human settlers and a variety of other life forms. The planet was first seen in the original 1977 film Star Wars, and has to date featured in a total five Star Wars theatrical films.\nIt is noted as the homeworld of the protagonist of the Star Wars saga, Luke Skywalker, and also of his father, Anakin Skywalker. Shots of the binary sunset over the Tatooine desert are considered to be an iconic image of the film series.[1][2]	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
2	Alderaan	24	364	12500	{TEMPERATE}	{GRASSLANDS,MOUNTAINS}	1	40	2000000000	Alderaan is a fictional planet featured in the Star Wars space opera franchise. It is blue-green in appearance and is depicted as a terrestrial planet with humanoid inhabitants. It is the home of Princess Leia, one of the lead characters in the film series. In the original 1977 film, Star Wars, Alderaan is blown up by the Death Star, a giant space station capable of destroying entire planets.[2]	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
3	Yavin IV	24	4818	10200	{TEMPERATE,TROPICAL}	{JUNGLES,RAINFORESTS}	1	8	1000	Yavin (also known as Yavin Prime) is a fictional planet in the Star Wars universe. It first appeared in the 1977 film Star Wars and is depicted as a large red gas giant with an extensive satellite system of moons. Within the Star Wars narrative, Yavin is noted as the hidden military base of the Rebel Alliance located on its fourth moon, known as Yavin IV.\nThe climactic space battle at the end of the film, in which the Rebel Alliance destroys the Death Star, takes place in orbit around the planet Yavin. In Star Wars fandom and the Star Wars expanded universe, this event is especially significant as it is used to mark an epoch in the fictional Star Wars universe. Events in Star Wars stories are typically dated in terms of years BBY (Before the Battle of Yavin) or ABY (After the Battle of Yavin).[1][2]	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
4	Hoth	23	549	7200	{FROZEN}	{TUNDRA,ICE_CAVES,MOUNTAINS}	1.10000000000000009	100	\N	Hoth is an ice planet in the Star Wars fictional universe. It first appeared in the 1980 film The Empire Strikes Back and has also been a setting in Star Wars books and video games.	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
5	Dagobah	23	341	8900	{MURKY}	{SWAMPS,JUNGLES}	\N	8	\N	Dagobah is a star system in the Star Wars films The Empire Strikes Back and Return of the Jedi. It also appears in a deleted scene from Revenge of the Sith. Yoda went into exile on Dagobah after a lightsaber battle with Darth Sidious.\nThe planet shown in Dagobah, in the Sluis sector, is a world of murky swamps, steaming bayous, and petrified forests, resembling Earth during the Carboniferous Period.	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
6	Bespin	12	5110	118000	{TEMPERATE}	{GAS_GIANT}	1.5	0	6000000	Bespin is a fictional planet, a gas giant in Star Wars films and books. The planet was first seen in the 1980 feature film The Empire Strikes Back. Since its introduction, Bespin has gained more specific characteristics in the Star Wars expanded universe.\nIn The Empire Strikes Back, Bespin's floating city Cloud City hovers suspended by an anti-gravity pod.	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
7	Endor	18	402	4900	{TEMPERATE}	{FORESTS,MOUNTAINS,LAKES}	0.849999999999999978	8	30000000	Endor is a planet in the Star Wars universe best known for its moon, known as the sanctuary moon, a forested world (moon) covered by giant trees.	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
8	Naboo	26	312	12120	{TEMPERATE}	{GRASSY_HILLS,SWAMPS,FORESTS,MOUNTAINS}	1	12	4500000000	Naboo is a planet in the fictional Star Wars universe with a mostly green terrain and which is the homeworld of two societies: the Gungans who dwell in underwater cities and the humans who live in colonies on the surface. Humans of Naboo have an electoral monarchy and maintain a peaceful culture that defends education, the arts, environmental protection and scientific achievements. The main capital of Naboo is Theed.\nLocated in the Chommell sector, Naboo is the home planet of Padmé Amidala and Jar Jar Binks, as well as Senator (later Supreme Chancellor and then Emperor) Palpatine. In Star Wars: Episode I – The Phantom Menace, it was the site of a blockade by the Trade Federation and the Battle of Naboo between the Federation and the native inhabitants. Naboo is seen in four films in the Star Wars series, having a prominent role in the first two prequels and glimpsed briefly in Revenge of the Sith and the 2004 DVD release of Return of the Jedi.\nTheed's architecture, while referencing Ancient Rome and other classical traditions, was heavily inspired by the Frank Lloyd Wright-designed Marin County Civic Center in California. (Skywalker Ranch and Industrial Light & Magic are both based in Marin County.)	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
9	Coruscant	24	368	12240	{TEMPERATE}	{CITYSCAPE,MOUNTAINS}	1	\N	1000000000000	Coruscant /ˈkɒrəsɑːnt/[2] is a planet in the fictional Star Wars universe. It first appeared onscreen in the 1997 Special Edition of Return of the Jedi, but was first mentioned in Timothy Zahn's 1991 novel Heir to the Empire. A city occupying an entire planet, it was renamed Imperial Center during the reign of the Galactic Empire (as depicted in the original films) and Yuuzhan'tar during the Yuuzhan Vong invasion (as depicted in the New Jedi Order novel series). The demonym and adjective form of the planet name is Coruscanti.\nCoruscant is, at various times, the capital of the Old Republic, the Galactic Empire, the New Republic, the Yuuzhan Vong Empire and the Galactic Alliance. Not only is Coruscant central to all these governing bodies, it is the navigational center of the galaxy, given that its hyperspace coordinates are (0,0,0). Due to its location and large population, roughly 2 trillion sentients, the galaxy's main trade routes — Perlemian Trade Route, Hydian Way, Corellian Run and Corellian Trade Spine — go through Coruscant, making it the richest and most influential world in the Star Wars galaxy. Coruscant is the sixth planet out of 11 planets in the Coruscant solar system, and has four moons; Centax-1, Centax-2, Centax-3, and Hesperidium.\nThe Galactic Standard Calendar was the standard measurement of time in the Star Wars galaxy. It centered on the Coruscant tropical year. The Coruscant solar cycle was 368 days long; with a day consisting of 360 NET degrees (or 24 standard hours).[3] Numerous epochs were used to determine calendar eras. The most recent of these calendar eras used the Battle of Yavin (i.e. the destruction of the first Death Star) as its epoch, or "year zero": BBY (Before the Battle of Yavin), and ABY (After the Battle of Yavin).[4] The earliest date in the Star Wars expanded universe as a whole is 13,000,000,000 BBY, which serves as the year the universe was created.	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
10	Kamino	27	463	19720	{TEMPERATE}	{OCEAN}	1	100	1000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
11	Geonosis	30	256	11370	{TEMPERATE,ARID}	{ROCK,DESERTS,MOUNTAINS,BARREN}	0.900000000000000022	5	100000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
12	Utapau	27	351	12900	{TEMPERATE,ARID,WINDY}	{SCRUBLANDS,SAVANNA,CANYONS,SINKHOLES}	1	0	95000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
13	Mustafar	36	412	4200	{HOT}	{VOLCANOES,LAVA_RIVERS,MOUNTAINS,CAVES}	1	0	20000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
14	Kashyyyk	26	381	12765	{TROPICAL}	{JUNGLES,FORESTS,LAKES,RIVERS}	1	60	45000000	Kashyyyk (/kəˈʃiːk/ ka-SHEEK, /ˈkæʃɪk/ KASH-ik or /ˌkæˈʃiː.aɪk/ ka-SHEE-ike), also known as Wookiee Planet C, is a fictional planet in the Star Wars universe. It is the tropical, forested home world of the Wookiees. According to interviews given by Star Wars creator George Lucas, the home of the Wookiees was originally intended to be the forest moon of Endor which plays a key role in the plot of the sixth film of the series, Return of the Jedi.[1] However, Lucas decided that since the Wookiee Chewbacca was clearly proficient with advanced technology (i.e. he was pilot and mechanic of the spaceship the Millennium Falcon and he repaired the damaged droid C-3PO), it would be confusing to show the Wookiees with a primitive, "stone age" culture on Endor. The Ewoks were created instead to populate the moon and to help fight the Imperial garrison stationed there.[2] Kashyyyk made appearances in the Star Wars Holiday Special and Star Wars: Episode III – Revenge of the Sith.	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
15	Polis Massa	24	590	0	{"ARTIFICIAL TEMPERATE"}	{AIRLESS_ASTEROID}	0.560000000000000053	0	1000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
16	Mygeeto	12	167	10088	{FRIGID}	{GLACIERS,MOUNTAINS,ICE_CANYONS}	1	\N	19000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
17	Felucia	34	231	9100	{HOT,HUMID}	{"FUNGUS FORESTS"}	0.75	\N	8500000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
18	Cato Neimoidia	25	278	0	{TEMPERATE,MOIST}	{MOUNTAINS,PLAINS,FORESTS,ROCK_ARCHES}	1	\N	10000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
19	Saleucami	26	392	14920	{HOT}	{CAVES,DESERTS,MOUNTAINS,VOLCANOES}	\N	\N	1400000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
20	Stewjon	\N	\N	0	{TEMPERATE}	{GRASS}	1	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
21	Eriadu	24	360	13490	{POLLUTED}	{CITYSCAPE}	1	\N	22000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
22	Corellia	25	329	11000	{TEMPERATE}	{PLAINS,URBAN,HILLS,FORESTS}	1	70	3000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
23	Rodia	29	305	7549	{HOT}	{JUNGLES,OCEANS,URBAN,SWAMPS}	1	60	1300000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
24	Nal Hutta	87	413	12150	{TEMPERATE}	{URBAN,OCEANS,SWAMPS,BOGS}	1	\N	7000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
25	Dantooine	25	378	9830	{TEMPERATE}	{OCEANS,SAVANNAS,MOUNTAINS,GRASSLANDS}	1	\N	1000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
26	Bestine IV	26	680	6400	{TEMPERATE}	{ISLANDS,OCEANS}	\N	98	62000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
27	Ord Mantell	26	334	14050	{TEMPERATE}	{PLAINS,SEAS,MESAS}	1	10	4000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
29	Trandosha	25	371	0	{ARID}	{MOUNTAINS,SEAS,GRASSLANDS,DESERTS}	0.619999999999999996	\N	42000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
30	Socorro	20	326	0	{ARID}	{DESERTS,MOUNTAINS}	1	\N	300000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
31	Mon Cala	21	398	11030	{TEMPERATE}	{OCEANS,REEFS,ISLANDS}	1	100	27000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
32	Chandrila	20	368	13500	{TEMPERATE}	{PLAINS,FORESTS}	1	40	1200000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
33	Sullust	20	263	12780	{SUPERHEATED}	{MOUNTAINS,VOLCANOES,ROCKY_DESERTS}	1	5	18500000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
34	Toydaria	21	184	7900	{TEMPERATE}	{SWAMPS,LAKES}	1	\N	11000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
35	Malastare	26	201	18880	{ARID,TEMPERATE,TROPICAL}	{SWAMPS,DESERTS,JUNGLES,MOUNTAINS}	1.56000000000000005	\N	2000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
36	Dathomir	24	491	10480	{TEMPERATE}	{FORESTS,DESERTS,SAVANNAS}	0.900000000000000022	\N	5200	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
37	Ryloth	30	305	10600	{TEMPERATE,ARID,SUBARTIC}	{MOUNTAINS,VALLEYS,DESERTS,TUNDRA}	1	5	1500000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
38	Aleen Minor	\N	\N	\N	\N	\N	\N	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
39	Vulpter	22	391	14900	{TEMPERATE,ARTIC}	{URBAN,BARREN}	1	\N	421000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
40	Troiken	\N	\N	\N	\N	{DESERTS,TUNDRA,RAINFORESTS,MOUNTAINS}	\N	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
41	Tund	48	1770	12190	\N	{BARREN,ASH}	\N	\N	0	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
42	Haruun Kal	25	383	10120	{TEMPERATE}	{TOXIC_CLOUDSEA,PLATEAUS,VOLCANOES}	0.979999999999999982	\N	705300	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
43	Cerea	27	386	\N	{TEMPERATE}	{VERDANT}	1	20	450000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
44	Glee Anselm	33	206	15600	{TROPICAL,TEMPERATE}	{LAKES,ISLANDS,SWAMPS,SEAS}	1	80	500000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
45	Iridonia	29	413	\N	\N	{ROCKY_CANYONS,ACID_POOLS}	\N	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
46	Tholoth	\N	\N	\N	\N	\N	\N	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
47	Iktotch	22	481	\N	{ARID,ROCKY,WINDY}	{ROCKY}	1	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
48	Quermia	\N	\N	\N	\N	\N	\N	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
49	Dorin	22	409	13400	{TEMPERATE}	\N	1	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
50	Champala	27	318	\N	{TEMPERATE}	{OCEANS,RAINFORESTS,PLATEAUS}	1	\N	3500000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
51	Mirial	\N	\N	\N	\N	{DESERTS}	\N	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
52	Serenno	\N	\N	\N	\N	{RAINFORESTS,RIVERS,MOUNTAINS}	\N	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
53	Concord Dawn	\N	\N	\N	\N	{JUNGLES,FORESTS,DESERTS}	\N	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
54	Zolan	\N	\N	\N	\N	\N	\N	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
55	Ojom	\N	\N	\N	{FRIGID}	{OCEANS,GLACIERS}	\N	100	500000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
56	Skako	27	384	\N	{TEMPERATE}	{URBAN,VINES}	1	\N	500000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
57	Muunilinst	28	412	13800	{TEMPERATE}	{PLAINS,FORESTS,HILLS,MOUNTAINS}	1	25	5000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
58	Shili	\N	\N	\N	{TEMPERATE}	{CITIES,SAVANNAHS,SEAS,PLAINS}	1	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
59	Kalee	23	378	13850	{ARID,TEMPERATE,TROPICAL}	{RAINFORESTS,CLIFFS,CANYONS,SEAS}	1	\N	4000000000	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
60	Umbara	\N	\N	\N	\N	\N	\N	\N	\N	\N	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
61	Jakku	\N	\N	\N	\N	{DESERTS}	\N	\N	\N	Jakku is a fictional desert planet first featured in the 2015 Star Wars film The Force Awakens. Remote, lawless, and inhospitable, it is the homeworld of main character Rey, played by Daisy Ridley. The film focuses on two distinct localities, Tuanul Village and Niima Outpost, near a starship graveyard.\nThe planet is also depicted in the 2017 Chuck Wendig novel, Star Wars: Aftermath: Empire's End.	2019-03-21 19:14:07.095112	2019-03-21 19:14:07.095112
\.


--
-- Data for Name: species; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.species (id, name, classification, average_height, skin_colors, hair_colors, eye_colors, average_lifespan, home_world, spoken_language, description, created, updated) FROM stdin;
1	Human	mammal	180	{CAUCASIAN,BLACK,ASIAN,HISPANIC}	{BLONDE,BROWN,BLACK,RED}	{BROWN,BLUE,GREEN,HAZEL,GREY,AMBER}	120	9	Galactic Basic	Modern humans (Homo sapiens, primarily ssp. Homo sapiens sapiens) are the only extant members of the subtribe Hominina, a branch of the tribe Hominini belonging to the family of great apes. They are characterized by erect posture and bipedal locomotion; high manual dexterity and heavy tool use compared to other animals; and a general trend toward larger, more complex brains and societies.[3][4]\nEarly hominins—particularly the australopithecines, whose brains and anatomy are in many ways more similar to ancestral non-human apes—are less often referred to as "human" than hominins of the genus Homo.[5] Several of these hominins used fire, occupied much of Eurasia, and gave rise to anatomically modern Homo sapiens in Africa about 200,000 years ago.[6][7] They began to exhibit evidence of behavioral modernity around 50,000 years ago. In several waves of migration, anatomically modern humans ventured out of Africa and populated most of the world.[8]\nThe spread of humans and their large and increasing population has had a profound impact on large areas of the environment and millions of native species worldwide. Advantages that explain this evolutionary success include a relatively larger brain with a particularly well-developed neocortex, prefrontal cortex and temporal lobes, which enable high levels of abstract reasoning, language, problem solving, sociality, and culture through social learning. Humans use tools to a much higher degree than any other animal, are the only extant species known to build fires and cook their food, and are the only extant species to clothe themselves and create and use numerous other technologies and arts.\nHumans are uniquely adept at using systems of symbolic communication (such as language and art) for self-expression and the exchange of ideas, and for organizing themselves into purposeful groups. Humans create complex social structures composed of many cooperating and competing groups, from families and kinship networks to political states. Social interactions between humans have established an extremely wide variety of values,[9] social norms, and rituals, which together form the basis of human society. Curiosity and the human desire to understand and influence the environment and to explain and manipulate phenomena (or events) has provided the foundation for developing science, philosophy, mythology, religion, anthropology, and numerous other fields of knowledge.\nThough most of human existence has been sustained by hunting and gathering in band societies,[10] increasing numbers of human societies began to practice sedentary agriculture approximately some 10,000 years ago,[11] domesticating plants and animals, thus allowing for the growth of civilization. These human societies subsequently expanded in size, establishing various forms of government, religion, and culture around the world, unifying people within regions to form states and empires. The rapid advancement of scientific and medical understanding in the 19th and 20th centuries led to the development of fuel-driven technologies and increased lifespans, causing the human population to rise exponentially. Today the global human population is estimated by the United Nations to be near 7.6 billion.[12]	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
2	Droid	artificial	\N	\N	\N	\N	10000	\N	\N	A droid is a fictional robot possessing some degree of artificial intelligence in the Star Wars science fiction franchise. Coined by special effects artist John Stears, the term is a clipped form of "android",[1] a word originally reserved for robots designed to look and act like a human.[2] The word "droid" has been a registered trademark of Lucasfilm Ltd since 1977.[3][4][5][6]	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
3	Wookiee	mammal	210	{GRAY}	{BLACK,BROWN}	{BLUE,GREEN,YELLOW,BROWN,GOLDEN,RED}	400	14	Shyriiwook	Wookiees (/ˈwʊkiːz/) are a fictional species of intelligent bipeds from the planet Kashyyyk in the Star Wars universe. They are taller, stronger, and hairier than humans and most (if not all) other humanoid species. The most notable Wookiee is Chewbacca, the copilot of Han Solo, who first appeared in the 1977 film Star Wars Episode IV: A New Hope.	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
4	Rodian	sentient	170	{GREEN,BLUE}	\N	{BLACK}	\N	23	Galactic Basic	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
5	Hutt	gastropod	300	{GREEN,BROWN,TAN}	\N	{YELLOW,RED}	1000	24	Huttese	The Hutts are a fictional alien race in the Star Wars universe. They appear in The Phantom Menace, Return of the Jedi and The Clone Wars, as well as the special edition release of A New Hope. They also appear in various Star Wars games, including those based on the movies, and the Knights of the Old Republic series. None of these are very friendly and all are criminally involved.[1] In the comic book series Tales of the Jedi: Golden Age of the Sith and Tales of the Jedi: The Fall of the Sith Empire, however, there is a Hutt character named Aarrba who is sympathetic to the main characters, Gav and Jori Daragon.	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
6	Yoda's species	mammal	66	{GREEN,YELLOW}	{BROWN,WHITE}	{BROWN,GREEN,YELLOW}	900	\N	Galactic basic	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
7	Trandoshan	reptile	200	{BROWN,GREEN}	{NONE}	{YELLOW,ORANGE}	\N	29	Dosh	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
8	Mon Calamari	amphibian	160	{RED,BLUE,BROWN,MAGENTA}	{NONE}	{YELLOW}	\N	31	Mon Calamarian	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
9	Ewok	mammal	100	{BROWN}	{WHITE,BROWN,BLACK}	{ORANGE,BROWN}	\N	7	Ewokese	Ewoks are a fictional race of small, mammaloid bipeds that appear in the Star Wars universe. They are hunter-gatherers resembling teddy bears that inhabit the forest moon of Endor and live in various arboreal huts and other simple dwellings. They first appeared in the 1983 film Return of the Jedi and have since appeared in two made-for-television films, Caravan of Courage: An Ewok Adventure (1984) and Ewoks: The Battle for Endor (1985), as well as a short-lived animated series and several books and games.	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
10	Sullustan	mammal	180	{PALE}	{NONE}	{BLACK}	\N	33	Sullutese	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
11	Neimodian	\N	180	{GREY,GREEN}	{NONE}	{RED,PINK}	\N	18	Neimoidia	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
12	Gungan	amphibian	190	{BROWN,GREEN}	{NONE}	{ORANGE}	\N	8	Gungan basic	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
13	Toydarian	mammal	120	{BLUE,GREEN,GREY}	{NONE}	{YELLOW}	91	34	Toydarian	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
14	Dug	mammal	100	{BROWN,PURPLE,GREY,RED}	{NONE}	{YELLOW,BLUE}	\N	35	Dugese	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
15	Twi'lek	mammals	200	{ORANGE,YELLOW,BLUE,GREEN,PINK,PURPLE,TAN}	{NONE}	{BLUE,BROWN,ORANGE,PINK}	\N	37	Twi'leki	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
16	Aleena	reptile	80	{BLUE,GRAY}	{NONE}	\N	79	38	Aleena	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
17	Vulptereen	\N	100	{GREY}	{NONE}	{YELLOW}	\N	39	vulpterish	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
18	Xexto	\N	125	{GREY,YELLOW,PURPLE}	{NONE}	{BLACK}	\N	40	Xextese	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
19	Toong	\N	200	{GREY,GREEN,YELLOW}	{NONE}	{ORANGE}	\N	41	Tundan	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
20	Cerean	mammal	200	{"PALE PINK"}	{RED,BLOND,BLACK,WHITE}	{HAZEL}	\N	43	Cerean	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
21	Nautolan	amphibian	180	{GREEN,BLUE,BROWN,RED}	{NONE}	{BLACK}	70	44	Nautila	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
22	Zabrak	mammal	180	{PALE,BROWN,RED,ORANGE,YELLOW}	{BLACK}	{BROWN,ORANGE}	\N	45	Zabraki	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
23	Tholothian	mammal	\N	{DARK}	\N	{BLUE,INDIGO}	\N	46	\N	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
24	Iktotchi	\N	180	{PINK}	{NONE}	{ORANGE}	\N	47	Iktotchese	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
25	Quermian	mammal	240	{WHITE}	{NONE}	{YELLOW}	86	48	Quermian	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
26	Kel Dor	\N	180	{PEACH,ORANGE,RED}	{NONE}	{BLACK,SILVER}	70	49	Kel Dor	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
27	Chagrian	amphibian	190	{BLUE}	{NONE}	{BLUE}	\N	50	Chagria	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
28	Geonosian	insectoid	178	{GREEN,BROWN}	{NONE}	{GREEN,HAZEL}	\N	11	Geonosian	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
29	Mirialan	mammal	180	{YELLOW,GREEN}	{BLACK,BROWN}	{BLUE,GREEN,RED,YELLOW,BROWN,ORANGE}	\N	51	Mirialan	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
30	Clawdite	reptilian	180	{GREEN,YELLOW}	{NONE}	{YELLOW}	70	54	Clawdite	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
31	Besalisk	amphibian	178	{BROWN}	{NONE}	{YELLOW}	75	55	besalisk	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
32	Kaminoan	amphibian	220	{GREY,BLUE}	{NONE}	{BLACK}	80	10	Kaminoan	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
33	Skakoan	mammal	\N	{GREY,GREEN}	{NONE}	\N	\N	56	Skakoan	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
34	Muun	mammal	190	{GREY,WHITE}	{NONE}	{BLACK}	100	57	Muun	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
35	Togruta	mammal	180	{RED,WHITE,ORANGE,YELLOW,GREEN,BLUE}	{NONE}	{RED,ORANGE,YELLOW,GREEN,BLUE,BLACK}	94	58	Togruti	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
36	Kaleesh	reptile	170	{BROWN,ORANGE,TAN}	{NONE}	{YELLOW}	80	59	Kaleesh	\N	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
37	Pau'an	mammal	190	\N	\N	\N	700	12	Utapese	The Pau'ans were a species of gaunt humanoids with gray, furrowed skin, and sunken eyes. Their arms ended in long, pentadactyl hands with an opposable thumb and sharp nails. Being carnivores with a preference for raw meat, they also had sharp, jagged teeth. In addition to those intimidating features, Pau'ans were tall—Tion Medon, a Port Administrator during the Clone Wars, stood at about 2.06 meters. Their speech organ allowed them to speak Galactic Basic. They had hypersensitive hearing, which prompted them to wear special coverings to protect their fragile ears. Their life-span was measured in centuries, which earned them the nickname of "Ancients." The Pau'ans' native language was Utapese, however they were also capable of speaking Basic.	2019-03-21 19:14:34.344726	2019-03-21 19:14:34.344726
\.


--
-- Data for Name: vehicles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vehicles (id, name, model, manufacturer, cost_in_credits, size, max_atmospheric_speed, crew, passengers, cargo_capacity, consumables, hyperdrive_rating, mglt, starship_class, created, updated) FROM stdin;
2	CR90 corvette	CR90 corvette	Corellian Engineering Corporation	3500000	150	950	165	600	3000000	1 year	2	60	corvette	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
3	Star Destroyer	Imperial I-class Star Destroyer	Kuat Drive Yards	150000000	1	975	47060	0	36000000	2 years	2	60	Star Destroyer	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
4	Sand Crawler	Digger Crawler	Corellia Mining Corporation	150000	36	30	46	30	50000	2 mons	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
5	Sentinel-class landing craft	Sentinel-class landing craft	Sienar Fleet Systems, Cyngus Spaceworks	240000	38	1000	5	75	180000	1 mon	1	70	landing craft	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
6	T-16 skyhopper	T-16 skyhopper	Incom Corporation	14500	10	1200	1	1	50	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
7	X-34 landspeeder	X-34 landspeeder	SoroSuub Corporation	10550	3	250	1	1	5	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
8	TIE/LN starfighter	Twin Ion Engine/Ln Starfighter	Sienar Fleet Systems	\N	6	1200	1	0	65	2 days	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
9	Death Star	DS-1 Orbital Battle Station	Imperial Department of Military Research, Sienar Fleet Systems	1000000000000	120000	\N	342953	843342	1000000000000	3 years	4	10	Deep Space Mobile Battlestation	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
10	Millennium Falcon	YT-1300 light freighter	Corellian Engineering Corporation	100000	34	1050	4	6	100000	2 mons	0	75	Light freighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
11	Y-wing	BTL Y-wing	Koensayr Manufacturing	134999	14	1000	2	0	110	7 days	1	80	assault starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
12	X-wing	T-65 X-wing	Incom Corporation	149999	12	1050	1	0	110	7 days	1	100	Starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
13	TIE Advanced x1	Twin Ion Engine Advanced x1	Sienar Fleet Systems	\N	9	1200	1	0	150	5 days	1	105	Starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
14	Snowspeeder	t-47 airspeeder	Incom corporation	\N	4	650	2	0	10	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
15	Executor	Executor-class star dreadnought	Kuat Drive Yards, Fondor Shipyards	1143350000	19	\N	279144	38000	250000000	6 years	2	40	Star dreadnought	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
16	TIE bomber	TIE/sa bomber	Sienar Fleet Systems	\N	7	850	1	0	0	2 days	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
17	Rebel transport	GR-75 medium transport	Gallofree Yards, Inc.	\N	90	650	6	90	19000000	6 mons	4	20	Medium transport	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
18	AT-AT	All Terrain Armored Transport	Kuat Drive Yards, Imperial Department of Military Research	\N	20	60	5	40	1000	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
19	AT-ST	All Terrain Scout Transport	Kuat Drive Yards, Imperial Department of Military Research	\N	2	90	2	0	200	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
20	Storm IV Twin-Pod cloud car	Storm IV Twin-Pod	Bespin Motors	75000	7	1500	2	0	10	1 day	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
21	Slave 1	Firespray-31-class patrol and attack	Kuat Systems Engineering	\N	21	1000	1	6	70000	1 mon	3	70	Patrol craft	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
22	Imperial shuttle	Lambda-class T-4a shuttle	Sienar Fleet Systems	240000	20	850	6	20	80000	2 mons	1	50	Armed government transport	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
23	EF76 Nebulon-B escort frigate	EF76 Nebulon-B escort frigate	Kuat Drive Yards	8500000	300	800	854	75	6000000	2 years	2	40	Escort ship	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
24	Sail barge	Modified Luxury Sail Barge	Ubrikkian Industries Custom Vehicle Division	285000	30	100	26	500	2000000	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
25	Bantha-II cargo skiff	Bantha-II	Ubrikkian Industries	8000	9	250	5	16	135000	1 day	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
26	TIE/IN interceptor	Twin Ion Engine Interceptor	Sienar Fleet Systems	\N	9	1250	1	0	75	2 days	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
27	Calamari Cruiser	MC80 Liberty type Star Cruiser	Mon Calamari shipyards	104000000	1200	\N	5400	1200	\N	2 years	1	60	Star Cruiser	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
28	A-wing	RZ-1 A-wing Interceptor	Alliance Underground Engineering, Incom Corporation	175000	9	1300	1	0	40	7 days	1	120	Starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
29	B-wing	A/SF-01 B-wing starfighter	Slayn & Korpil	220000	16	950	1	0	45	7 days	2	91	Assault Starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
30	Imperial Speeder Bike	74-Z speeder bike	Aratech Repulsor Company	8000	3	360	1	1	4	1 day	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
31	Republic Cruiser	Consular-class cruiser	Corellian Engineering Corporation	\N	115	900	9	16	\N	\N	2	\N	Space cruiser	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
32	Droid control ship	Lucrehulk-class Droid Control Ship	Hoersch-Kessel Drive, Inc.	\N	3170	\N	175	139000	4000000000	500 days	2	\N	Droid control ship	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
33	Vulture Droid	Vulture-class droid starfighter	Haor Chall Engineering, Baktoid Armor Workshop	\N	3	1200	0	0	0	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
34	Multi-Troop Transport	Multi-Troop Transport	Baktoid Armor Workshop	138000	31	35	4	112	12000	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
35	Armored Assault Tank	Armoured Assault Tank	Baktoid Armor Workshop	\N	9	55	4	6	\N	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
36	Single Trooper Aerial Platform	Single Trooper Aerial Platform	Baktoid Armor Workshop	2500	2	400	1	0	0	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
37	C-9979 landing craft	C-9979 landing craft	Haor Chall Engineering	200000	210	587	140	284	1800000	1 day	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
38	Tribubble bongo	Tribubble bongo	Otoh Gunga Bongameken Cooperative	\N	15	85	1	2	1600	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
39	Naboo fighter	N-1 starfighter	Theed Palace Space Vessel Engineering Corps	200000	11	1100	1	0	65	7 days	1	\N	Starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
40	Naboo Royal Starship	J-type 327 Nubian royal starship	Theed Palace Space Vessel Engineering Corps, Nubia Star Drives	\N	76	920	8	\N	\N	\N	1	\N	yacht	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
41	Scimitar	Star Courier	Republic Sienar Systems	55000000	26	1180	1	6	2500000	30 days	1	\N	Space Transport	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
42	Sith speeder	FC-20 speeder bike	Razalon	4000	1	180	1	0	2	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
43	J-type diplomatic barge	J-type diplomatic barge	Theed Palace Space Vessel Engineering Corps, Nubia Star Drives	2000000	39	2000	5	10	\N	1 year	0	\N	Diplomatic barge	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
44	Zephyr-G swoop bike	Zephyr-G swoop bike	Mobquet Swoops and Speeders	5750	3	350	1	1	200	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
45	Koro-2 Exodrive airspeeder	Koro-2 Exodrive airspeeder	Desler Gizh Outworld Mobility Corporation	\N	6	800	1	1	80	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
46	XJ-6 airspeeder	XJ-6 airspeeder	Narglatch AirTech prefabricated kit	\N	6	720	1	1	\N	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
47	AA-9 Coruscant freighter	Botajef AA-9 Freighter-Liner	Botajef Shipyards	\N	390	\N	\N	30000	\N	\N	\N	\N	freighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
48	Jedi starfighter	Delta-7 Aethersprite-class interceptor	Kuat Systems Engineering	180000	8	1150	1	0	60	7 days	1	\N	Starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
49	H-type Nubian yacht	H-type Nubian yacht	Theed Palace Space Vessel Engineering Corps	\N	47	8000	4	\N	\N	\N	0	\N	yacht	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
50	LAAT/i	Low Altitude Assault Transport/infrantry	Rothana Heavy Engineering	\N	17	620	6	30	170	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
51	LAAT/c	Low Altitude Assault Transport/carrier	Rothana Heavy Engineering	\N	28	620	1	0	40000	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
52	Republic Assault ship	Acclamator I-class assault ship	Rothana Heavy Engineering	\N	752	\N	700	16000	11250000	2 years	0	\N	assault ship	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
53	AT-TE	All Terrain Tactical Enforcer	Rothana Heavy Engineering, Kuat Drive Yards	\N	13	60	6	36	10000	21 days	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
54	SPHA	Self-Propelled Heavy Artillery	Rothana Heavy Engineering	\N	140	35	25	30	500	7 days	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
55	Flitknot speeder	Flitknot speeder	Huppla Pasa Tisc Shipwrights Collective	8000	2	634	1	0	\N	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
56	Neimoidian shuttle	Sheathipede-class transport shuttle	Haor Chall Engineering	\N	20	880	2	6	1000	7 days	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
57	Geonosian starfighter	Nantex-class territorial defense	Huppla Pasa Tisc Shipwrights Collective	\N	9	20000	1	0	\N	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
58	Solar Sailer	Punworcca 116-class interstellar sloop	Huppla Pasa Tisc Shipwrights Collective	35700	15	1600	3	11	240	7 days	1	\N	yacht	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
59	Trade Federation cruiser	Providence-class carrier/destroyer	Rendili StarDrive, Free Dac Volunteers Engineering corps.	125000000	1088	1050	600	48247	50000000	4 years	1	\N	capital ship	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
60	Tsmeu-6 personal wheel bike	Tsmeu-6 personal wheel bike	Z-Gomot Ternbuell Guppat Corporation	15000	3	330	1	1	10	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
61	Theta-class T-2c shuttle	Theta-class T-2c shuttle	Cygnus Spaceworks	1000000	18	2000	5	16	50000	56 days	1	\N	transport	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
62	Emergency Firespeeder	Fire suppression speeder	\N	\N	\N	\N	2	\N	\N	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
63	Republic attack cruiser	Senator-class Star Destroyer	Kuat Drive Yards, Allanteen Six shipyards	59000000	1137	975	7400	2000	20000000	2 years	1	\N	star destroyer	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
64	Naboo star skiff	J-type star skiff	Theed Palace Space Vessel Engineering Corps/Nubia Star Drives, Incorporated	\N	29	1050	3	3	\N	\N	0	\N	yacht	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
65	Jedi Interceptor	Eta-2 Actis-class light interceptor	Kuat Systems Engineering	320000	5	1500	1	0	60	2 days	1	\N	starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
66	arc-170	Aggressive Reconnaissance-170 starfighte	Incom Corporation, Subpro Corporation	155000	14	1000	3	0	110	5 days	1	100	starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
67	Droid tri-fighter	tri-fighter	Colla Designs, Phlac-Arphocc Automata Industries	20000	5	1180	1	0	0	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
68	Banking clan frigte	Munificent-class star frigate	Hoersch-Kessel Drive, Inc, Gwori Revolutionary Industries	57000000	825	\N	200	\N	40000000	2 years	1	\N	cruiser	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
69	Oevvaor jet catamaran	Oevvaor jet catamaran	Appazanna Engineering Works	12125	15	420	2	2	50	3 days	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
70	Raddaugh Gnasp fluttercraft	Raddaugh Gnasp fluttercraft	Appazanna Engineering Works	14750	7	310	2	0	20	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
71	Clone turbo tank	HAVw A6 Juggernaut	Kuat Drive Yards	350000	49	160	20	300	30000	20 days	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
72	Corporate Alliance tank droid	NR-N99 Persuader-class droid enforcer	Techno Union	49000	10	100	0	4	0	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
73	Droid gunship	HMP droid gunship	Baktoid Fleet Ordnance, Haor Chall Engineering	60000	12	820	0	0	0	\N	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
74	Belbullab-22 starfighter	Belbullab-22 starfighter	Feethan Ottraw Scalable Assemblies	168000	6	1100	1	0	140	7 days	6	\N	starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
75	V-wing	Alpha-3 Nimbus-class V-wing starfighter	Kuat Systems Engineering	102500	7	1050	1	0	60	15:00:00	1	\N	starfighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
76	AT-RT	All Terrain Recon Transport	Kuat Drive Yards	40000	3	90	1	0	20	1 day	\N	\N	\N	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
77	T-70 X-wing fighter	T-70 X-wing fighter	Incom	\N	\N	\N	1	\N	\N	\N	\N	\N	fighter	2019-03-21 19:14:19.676908	2019-03-21 19:14:19.676908
\.


--
-- Name: characters_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.characters_id_seq', 1, false);


--
-- Name: constants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.constants_id_seq', 1, false);


--
-- Name: films_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.films_id_seq', 1, false);


--
-- Name: planets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.planets_id_seq', 1, false);


--
-- Name: species_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.species_id_seq', 1, false);


--
-- Name: vehicles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.vehicles_id_seq', 1, false);


--
-- Name: characters characters_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.characters
    ADD CONSTRAINT characters_pkey PRIMARY KEY (id);


--
-- Name: constants constants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.constants
    ADD CONSTRAINT constants_pkey PRIMARY KEY (id);


--
-- Name: constants constants_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.constants
    ADD CONSTRAINT constants_unique UNIQUE (constant_type, constant_value);


--
-- Name: films films_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.films
    ADD CONSTRAINT films_pkey PRIMARY KEY (id);


--
-- Name: planets planets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.planets
    ADD CONSTRAINT planets_pkey PRIMARY KEY (id);


--
-- Name: species species_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.species
    ADD CONSTRAINT species_pkey PRIMARY KEY (id);


--
-- Name: vehicles vehicles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT vehicles_pkey PRIMARY KEY (id);


--
-- Name: characters characters_insert_created; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER characters_insert_created BEFORE INSERT ON public.characters FOR EACH ROW EXECUTE PROCEDURE public.create_insert_created();


--
-- Name: characters characters_notify_deleted; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER characters_notify_deleted AFTER DELETE ON public.characters REFERENCING OLD TABLE AS old_table FOR EACH STATEMENT EXECUTE PROCEDURE public.notify_deleted('characters');


--
-- Name: characters characters_notify_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER characters_notify_updated AFTER UPDATE ON public.characters REFERENCING NEW TABLE AS new_table FOR EACH STATEMENT EXECUTE PROCEDURE public.notify_updated('characters');


--
-- Name: characters characters_update_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER characters_update_updated BEFORE UPDATE ON public.characters FOR EACH ROW EXECUTE PROCEDURE public.create_update_updated();


--
-- Name: constants constants_insert_created; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER constants_insert_created BEFORE INSERT ON public.constants FOR EACH ROW EXECUTE PROCEDURE public.create_insert_created();


--
-- Name: constants constants_update_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER constants_update_updated BEFORE UPDATE ON public.constants FOR EACH ROW EXECUTE PROCEDURE public.create_update_updated();


--
-- Name: films films_insert_created; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER films_insert_created BEFORE INSERT ON public.films FOR EACH ROW EXECUTE PROCEDURE public.create_insert_created();


--
-- Name: films films_notify_deleted; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER films_notify_deleted AFTER DELETE ON public.films REFERENCING OLD TABLE AS old_table FOR EACH STATEMENT EXECUTE PROCEDURE public.notify_deleted('films');


--
-- Name: films films_notify_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER films_notify_updated AFTER UPDATE ON public.films REFERENCING NEW TABLE AS new_table FOR EACH STATEMENT EXECUTE PROCEDURE public.notify_updated('films');


--
-- Name: films films_update_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER films_update_updated BEFORE UPDATE ON public.films FOR EACH ROW EXECUTE PROCEDURE public.create_update_updated();


--
-- Name: planets planets_insert_created; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER planets_insert_created BEFORE INSERT ON public.planets FOR EACH ROW EXECUTE PROCEDURE public.create_insert_created();


--
-- Name: planets planets_notify_deleted; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER planets_notify_deleted AFTER DELETE ON public.planets REFERENCING OLD TABLE AS old_table FOR EACH STATEMENT EXECUTE PROCEDURE public.notify_deleted('planets');


--
-- Name: planets planets_notify_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER planets_notify_updated AFTER UPDATE ON public.planets REFERENCING NEW TABLE AS new_table FOR EACH STATEMENT EXECUTE PROCEDURE public.notify_updated('planets');


--
-- Name: planets planets_update_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER planets_update_updated BEFORE UPDATE ON public.planets FOR EACH ROW EXECUTE PROCEDURE public.create_update_updated();


--
-- Name: species species_insert_created; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER species_insert_created BEFORE INSERT ON public.species FOR EACH ROW EXECUTE PROCEDURE public.create_insert_created();


--
-- Name: species species_notify_deleted; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER species_notify_deleted AFTER DELETE ON public.species REFERENCING OLD TABLE AS old_table FOR EACH STATEMENT EXECUTE PROCEDURE public.notify_deleted('species');


--
-- Name: species species_notify_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER species_notify_updated AFTER UPDATE ON public.species REFERENCING NEW TABLE AS new_table FOR EACH STATEMENT EXECUTE PROCEDURE public.notify_updated('species');


--
-- Name: species species_update_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER species_update_updated BEFORE UPDATE ON public.species FOR EACH ROW EXECUTE PROCEDURE public.create_update_updated();


--
-- Name: vehicles vehicles_insert_created; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER vehicles_insert_created BEFORE INSERT ON public.vehicles FOR EACH ROW EXECUTE PROCEDURE public.create_insert_created();


--
-- Name: vehicles vehicles_notify_deleted; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER vehicles_notify_deleted AFTER DELETE ON public.vehicles REFERENCING OLD TABLE AS old_table FOR EACH STATEMENT EXECUTE PROCEDURE public.notify_deleted('vehicles');


--
-- Name: vehicles vehicles_notify_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER vehicles_notify_updated AFTER UPDATE ON public.vehicles REFERENCING NEW TABLE AS new_table FOR EACH STATEMENT EXECUTE PROCEDURE public.notify_updated('vehicles');


--
-- Name: vehicles vehicles_update_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER vehicles_update_updated BEFORE UPDATE ON public.vehicles FOR EACH ROW EXECUTE PROCEDURE public.create_update_updated();


--
-- Name: characters character_home_world; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.characters
    ADD CONSTRAINT character_home_world FOREIGN KEY (home_world) REFERENCES public.planets(id);


--
-- Name: characters character_species; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.characters
    ADD CONSTRAINT character_species FOREIGN KEY (species) REFERENCES public.species(id);


--
-- Name: species species_home_world; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.species
    ADD CONSTRAINT species_home_world FOREIGN KEY (home_world) REFERENCES public.planets(id);


--
-- PostgreSQL database dump complete
--

