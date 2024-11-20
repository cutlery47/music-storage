\connect music 


-- Таблица для хранения песен
CREATE TABLE music_schema.songs(
    id              music_schema.uuid_key       PRIMARY KEY,
    group_name      music_schema.string,
    song_name       music_schema.string,

    UNIQUE(group_name, song_name)
);

-- Таблица для хранения информации о песнях
CREATE TABLE music_schema.songs_details(
    id              music_schema.uuid_key       PRIMARY KEY,
    song_id         UUID                        REFERENCES music_schema.songs(id) ON DELETE CASCADE,
    released_at     date,
    link            music_schema.string,

    UNIQUE(song_id)
);

-- Таблица для хранения куплетов песен
CREATE TABLE music_schema.songs_verses(
    id              music_schema.uuid_key        PRIMARY KEY,
    song_id         UUID                         REFERENCES music_schema.songs(id) ON DELETE CASCADE,
    verse_id        music_schema.pos_int,
    verse           text,

    UNIQUE(song_id, verse_id)
)