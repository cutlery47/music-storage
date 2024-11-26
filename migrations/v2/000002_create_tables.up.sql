-- Таблица для хранения групп
CREATE TABLE IF NOT EXISTS music_schema.groups(
    id               music_schema.uuid_key       PRIMARY KEY,
    group_name       music_schema.string,

    UNIQUE(group_name)
);

-- Таблица для хранения песен
CREATE TABLE IF NOT EXISTS music_schema.songs(
    id          music_schema.uuid_key       PRIMARY KEY,
    group_id    UUID                        REFERENCES music_schema.groups(id) ON DELETE CASCADE,
    song_name   music_schema.string         
);

-- Таблица для хранения куплетов песен
CREATE TABLE IF NOT EXISTS music_schema.songs_verses(
    id          music_schema.uuid_key       PRIMARY KEY,
    song_id     UUID                        REFERENCES music_schema.songs(id) ON DELETE CASCADE,
    verse_num   music_schema.pos_int,
    verse       text,

    UNIQUE(song_id, verse_num)
);

-- Таблица для хранения информации о песнях
CREATE TABLE IF NOT EXISTS music_schema.songs_details(
    id              music_schema.uuid_key       PRIMARY KEY,
    song_id         UUID                        REFERENCES music_schema.songs(id) ON DELETE CASCADE,
    released_at     date,
    link            music_schema.string,

    UNIQUE(song_id)
);