package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cutlery47/music-storage/internal/models"
	"github.com/google/uuid"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Repository interface {
	// Добавление информации о песне
	Create(ctx context.Context, song models.SongWithDetailSplit) error
	// Получение информации о песнях по произвольным фильтрам
	Read(ctx context.Context, limit, offset int, filter models.Filter) ([]models.SongWithDetail, error)
	// Получение текста песни по куплетам
	ReadText(ctx context.Context, limit, offset int, song models.Song) ([]string, error)
	// Получение информации о конкретной песне
	ReadDetail(ctx context.Context, song models.Song) (models.SongDetail, error)
	// Обновление информации о песне
	Update(ctx context.Context, song models.Song, upd models.SongWithDetailSplit) error
	// Удаление песни
	Delete(ctx context.Context, song models.Song) error
}

type MusicRepository struct {
	db *sql.DB
}

func NewMusicRepository(url string) (*MusicRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &MusicRepository{
		db: db,
	}, nil
}

func (mr *MusicRepository) Create(ctx context.Context, song models.SongWithDetailSplit) error {
	queryInsertSong :=
		`
	INSERT INTO music_schema.songs
	(group_name, song_name)
	VALUES 
	($1, $2)
	RETURNING id
	`

	queryInsertDetail :=
		`
	INSERT INTO music_schema.songs_details
	(song_id, released_at, link)
	VALUES
	($1, $2, $3)
	`

	queryInsertText :=
		`
	INSERT INTO music_schema.songs_verses
	(song_id, verse_id, verse)
	VALUES
	`

	for i := range song.Verses {
		queryInsertText += fmt.Sprintf("($%v, $%v, $%v),\n", (i+1)*3-2, (i+1)*3-1, (i+1)*3)
	}
	queryInsertText = strings.TrimSuffix(queryInsertText, ",\n")

	tx, err := mr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("mr.db.BeginTx: %v", err)
	}
	defer tx.Rollback()

	stmtInsertSong, err := tx.PrepareContext(ctx, queryInsertSong)
	if err != nil {
		return fmt.Errorf("tx.PrepareContext: %v", err)
	}

	stmtInsertDetail, err := tx.PrepareContext(ctx, queryInsertDetail)
	if err != nil {
		return fmt.Errorf("tx.PrepareContext: %v", err)
	}

	stmtInsertText, err := tx.PrepareContext(ctx, queryInsertText)
	if err != nil {
		return fmt.Errorf("tx.PrepareContext: %v", err)
	}

	var id uuid.UUID

	res := stmtInsertSong.QueryRowContext(ctx, song.GroupName, song.SongName)
	if err := res.Scan(&id); err != nil {
		// проверка на уникальность
		if pqerr, ok := err.(*pq.Error); ok && pqerr.Code == "23505" {
			return ErrAlreadyExists
		}
		return fmt.Errorf("res.Scan: %v", err)
	}

	_, err = stmtInsertDetail.ExecContext(ctx, id, song.ReleaseDate, song.Link)
	if err != nil {
		return fmt.Errorf("stmtInsertDetail.ExecContext: %v", err)
	}

	var vals []interface{}
	for i, verse := range song.Verses {
		vals = append(vals, id, i+1, verse)
	}

	_, err = stmtInsertText.ExecContext(ctx, vals...)
	if err != nil {
		return fmt.Errorf("stmtInsertText.ExecContext")
	}

	return tx.Commit()
}

func (mr *MusicRepository) Read(ctx context.Context, limit, offset int, filter models.Filter) ([]models.SongWithDetail, error) {
	var appliedFilters []interface{}

	query :=
		`
	SELECT s.group_name, s.song_name, sd.released_at, sd.link
	FROM
	music_schema.songs AS s
	JOIN
	music_schema.songs_details AS sd
	ON s.id = sd.song_id
	WHERE
	`

	query = mr.applyFilters(query, filter, limit, offset, &appliedFilters)

	stmt, err := mr.db.PrepareContext(ctx, query)
	if err != nil {
		return []models.SongWithDetail{}, fmt.Errorf("mr.db.PrepareContext: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, appliedFilters...)
	if err != nil {
		return []models.SongWithDetail{}, fmt.Errorf("stmt.QueryContext: %v", err)
	}

	songs := []models.SongWithDetail{}
	for rows.Next() {
		song := models.SongWithDetail{}
		rows.Scan(&song.GroupName, &song.SongName, &song.ReleaseDate, &song.Link)
		songs = append(songs, song)
	}

	return songs, nil
}

func (mr *MusicRepository) ReadDetail(ctx context.Context, song models.Song) (models.SongDetail, error) {
	query :=
		`
	SELECT sd.released_at, sd.link
	FROM 
	music_schema.songs AS s
	JOIN
	music_schema.songs_details AS sd
	ON
	s.id = sd.song_id
	WHERE
	s.group_name = $1 AND s.song_name = $2
	`

	stmt, err := mr.db.PrepareContext(ctx, query)
	if err != nil {
		return models.SongDetail{}, fmt.Errorf("mr.db.PrepareContext: %v", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, song.GroupName, song.SongName)

	detail := models.SongDetail{}
	if err := row.Scan(&detail.ReleaseDate, &detail.Link); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.SongDetail{}, ErrNotFound
		}
		return models.SongDetail{}, fmt.Errorf("row.Scan: %v", err)
	}

	return detail, nil
}

func (mr *MusicRepository) ReadText(ctx context.Context, limit, offset int, song models.Song) ([]string, error) {
	query :=
		`
	SELECT sv.verse
	FROM 
	music_schema.songs AS s
	JOIN 
	music_schema.songs_verses AS sv
	ON 
	s.id = sv.song_id
	WHERE 
	s.group_name = $1 AND s.song_name = $2
	ORDER BY sv.verse_id
	LIMIT $3
	OFFSET $4 
	`

	stmt, err := mr.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("mr.db.PrepareContext: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, song.GroupName, song.SongName, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("stmt.ExecContext: %v", err)
	}

	verses := []string{}
	for rows.Next() {
		verse := ""
		rows.Scan(&verse)
		verses = append(verses, verse)
	}

	return verses, nil
}

func (mr *MusicRepository) Update(ctx context.Context, song models.Song, upd models.SongWithDetailSplit) error {
	return nil
}

func (mr *MusicRepository) Delete(ctx context.Context, song models.Song) error {
	query :=
		`
	DELETE FROM music_schema.songs AS s
	WHERE s.group_name = $1 AND s.song_name = $2;
	`

	stmt, err := mr.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("mr.db.PrepareContext: %v", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, song.GroupName, song.SongName)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext: %v", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("res.RowsAffected: %v", err)
	}

	log.Println("delete rows affected:", affected)

	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

// Принимаем структуру, содержащую всевозможные фильтры для поиска песни, а также лимит и оффсет для пагинации.
// Слайс applied хранит значения фильтров, эти значения затем передаются в качестве аргументов prepared statement.
func (mr *MusicRepository) applyFilters(query string, filter models.Filter, limit, offset int, applied *[]any) string {
	filterCount := 0

	if filter.Group != nil {
		filterCount++
		query += fmt.Sprintf("group_name = $%v\n", filterCount)
		*applied = append(*applied, *filter.Group)
	}

	if filter.Song != nil {
		if filterCount > 0 {
			query += "AND\n"
		}
		filterCount++
		query += fmt.Sprintf("song_name = $%v\n", filterCount)
		*applied = append(*applied, *filter.Song)
	}

	if filter.ReleasedAfter != nil {
		if filterCount > 0 {
			query += "AND\n"
		}
		filterCount++
		query += fmt.Sprintf("released_at >= $%v\n", filterCount)
		*applied = append(*applied, *filter.ReleasedAfter)
	}

	if filter.ReleasedBefore != nil {
		if filterCount > 0 {
			query += "AND\n"
		}
		filterCount++
		query += fmt.Sprintf("released_at <= $%v\n", filterCount)
		*applied = append(*applied, *filter.ReleasedBefore)
	}

	if filterCount == 0 {
		query = strings.TrimSuffix(query, "WHERE\n\t")
	}

	filterCount++
	query += fmt.Sprintf("LIMIT $%v OFFSET $%v;", filterCount, filterCount+1)
	*applied = append(*applied, limit, offset)

	return query
}
