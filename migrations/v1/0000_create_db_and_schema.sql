CREATE DATABASE music;

\connect music

CREATE SCHEMA music_schema;

CREATE DOMAIN music_schema.uuid_key as UUID
DEFAULT gen_random_uuid()
NOT NULL;

CREATE DOMAIN music_schema.string as VARCHAR(256)
NOT NULL;

CREATE DOMAIN music_schema.pos_int as INTEGER
CHECK (
    VALUE > 0
);
