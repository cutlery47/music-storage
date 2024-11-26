\connect music 

INSERT INTO music_schema.songs (group_name, song_name)
VALUES 
    ('The Beatles', 'Hey Jude'),
    ('Queen', 'Bohemian Rhapsody'),
    ('Nirvana', 'Smells Like Teen Spirit'),
    ('Linkin Park', 'In the End'),
    ('Coldplay', 'Fix You');


WITH song_ids AS (
    SELECT id, song_name FROM music_schema.songs
)

INSERT INTO music_schema.songs_details (song_id, released_at, link)
VALUES 
    ((SELECT id FROM song_ids WHERE song_name = 'Hey Jude'), '1968-08-26', 'https://example.com/hey_jude'),
    ((SELECT id FROM song_ids WHERE song_name = 'Bohemian Rhapsody'), '1975-10-31', 'https://example.com/bohemian_rhapsody'),
    ((SELECT id FROM song_ids WHERE song_name = 'Smells Like Teen Spirit'), '1991-09-10', 'https://example.com/smells_like_teen_spirit'),
    ((SELECT id FROM song_ids WHERE song_name = 'In the End'), '2001-10-09', 'https://example.com/in_the_end'),
    ((SELECT id FROM song_ids WHERE song_name = 'Fix You'), '2005-09-05', 'https://example.com/fix_you');

WITH song_ids AS (
    SELECT id, song_name FROM music_schema.songs
)

INSERT INTO music_schema.songs_verses (song_id, verse_id, verse)
VALUES
    ((SELECT id FROM song_ids WHERE song_name = 'Hey Jude'), 1, 'Hey Jude, dont make it bad'),
    ((SELECT id FROM song_ids WHERE song_name = 'Hey Jude'), 2, 'Take a sad song and make it better'),

    ((SELECT id FROM song_ids WHERE song_name = 'Bohemian Rhapsody'), 1, 'Is this the real life? Is this just fantasy?'),
    ((SELECT id FROM song_ids WHERE song_name = 'Bohemian Rhapsody'), 2, 'Caught in a landslide, no escape from reality'),

    ((SELECT id FROM song_ids WHERE song_name = 'Smells Like Teen Spirit'), 1, 'Load up on guns, bring your friends'),
    ((SELECT id FROM song_ids WHERE song_name = 'Smells Like Teen Spirit'), 2, 'Its fun to lose and to pretend'),

    ((SELECT id FROM song_ids WHERE song_name = 'In the End'), 1, 'It starts with one thing'),
    ((SELECT id FROM song_ids WHERE song_name = 'In the End'), 2, 'I dont know why, it doesnt even matter how hard you try'),

    ((SELECT id FROM song_ids WHERE song_name = 'Fix You'), 1, 'When you try your best but you dont succeed'),
    ((SELECT id FROM song_ids WHERE song_name = 'Fix You'), 2, 'When you get what you want but not what you need');