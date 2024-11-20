package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/cutlery47/music-storage/internal/models"

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
	Update(ctx context.Context, song models.Song, upd models.SongSplitText) error
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
	return nil
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

	row := stmt.QueryRowContext(ctx, song.GroupName, song.SongName)

	detail := models.SongDetail{}
	if err := row.Scan(&detail.ReleaseDate, &detail.Link); err != nil {
		return models.SongDetail{}, fmt.Errorf("row.Scan: %v", err)
	}

	return detail, nil
}

func (mr *MusicRepository) ReadText(ctx context.Context, limit, offset int, song models.Song) ([]string, error) {
	return []string{}, nil
}

func (mr *MusicRepository) Update(ctx context.Context, song models.Song, upd models.SongSplitText) error {
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
		return ErrSongNotFound
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
