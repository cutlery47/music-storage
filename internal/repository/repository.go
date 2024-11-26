package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cutlery47/music-storage/internal/config"
	"github.com/cutlery47/music-storage/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

// Repository impl
type MusicRepository struct {
	db *sql.DB
}

func NewMusicRepository(ctx context.Context, conf config.PostgresConfig) (*MusicRepository, error) {
	url := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=%v",
		conf.PostgresUser,
		conf.PostgresPassword,
		conf.PostgresHost,
		conf.PostgresPort,
		conf.PostgresDB,
		conf.PostgresSSL,
	)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	toCtx, cancel := context.WithTimeout(ctx, conf.PostgresTimeout*10)
	defer cancel()

	if err := db.PingContext(toCtx); err != nil {
		return nil, fmt.Errorf("couldn't establish connection with postgres: %v", err)
	}
	logrus.Debug("sucessfully established postgres connection!")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres.WithInstance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%v", conf.PostgresMigrations), conf.PostgresDB, driver)
	if err != nil {
		return nil, fmt.Errorf("migrate.New: %v", err)
	}

	logrus.Debug("applying migrations...")
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Debug("nothing to migrate")
		} else {
			return nil, fmt.Errorf("error when migrating: %v", err)
		}
	} else {
		logrus.Debug("migrated successfully!")
	}
	defer m.Close()

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

	// добавляем метки о всех куплетах, что хотим добавить
	for i := range song.Verses {
		queryInsertText += fmt.Sprintf("($%v, $%v, $%v),\n", (i+1)*3-2, (i+1)*3-1, (i+1)*3)
	}
	queryInsertText = strings.TrimSuffix(queryInsertText, ",\n")

	tx, err := mr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("mr.db.BeginTx: %v", err)
	}
	defer tx.Rollback()

	var id uuid.UUID

	res := tx.QueryRowContext(ctx, queryInsertSong, song.GroupName, song.SongName)
	if err := res.Scan(&id); err != nil {
		// проверка на уникальность
		if pqerr, ok := err.(*pq.Error); ok && pqerr.Code == "23505" {
			return ErrAlreadyExists
		}
		return fmt.Errorf("res.Scan: %v", err)
	}

	_, err = tx.ExecContext(ctx, queryInsertDetail, id, song.ReleaseDate, song.Link)
	if err != nil {
		return fmt.Errorf("st.ExecContext: %v", err)
	}

	// собираем значения для prepared statements
	var vals []interface{}
	for i, verse := range song.Verses {
		vals = append(vals, id, i+1, verse)
	}

	_, err = tx.ExecContext(ctx, queryInsertText, vals...)
	if err != nil {
		return fmt.Errorf("tx.ExecContext")
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

	rows, err := mr.db.QueryContext(ctx, query, appliedFilters...)
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

	row := mr.db.QueryRowContext(ctx, query, song.GroupName, song.SongName)

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

	rows, err := mr.db.QueryContext(ctx, query, song.GroupName, song.SongName, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("stmt.ExecContext: %v", err)
	}

	verses := []string{}
	for rows.Next() {
		verse := ""
		rows.Scan(&verse)
		verses = append(verses, verse)
	}

	if len(verses) == 0 {
		return nil, ErrNotFound
	}

	return verses, nil
}

func (mr *MusicRepository) Update(ctx context.Context, song models.Song, upd models.SongWithDetailSplit) error {
	queryUpdateSong :=
		`
	UPDATE music_schema.songs
	SET 
	group_name = $1,
	song_name = $2
	WHERE 
	group_name = $3 AND song_name = $4
	RETURNING id
	`

	queryUpdateDetail :=
		`
	UPDATE music_schema.songs_details
	SET 
	released_at = $1,
	link = $2
	WHERE 
	song_id = $3
	`

	queryDeleteOldVerses :=
		`
	DELETE FROM music_schema.songs_verses AS sv
	WHERE sv.song_id = $1
	`

	queryAddNewVerses :=
		`
	INSERT INTO music_schema.songs_verses AS sv
	(song_id, verse_id, verse)
	VALUES
	`

	for i := range upd.Verses {
		queryAddNewVerses += fmt.Sprintf("($%v, $%v, $%v),\n", (i+1)*3-2, (i+1)*3-1, (i+1)*3)
	}
	queryAddNewVerses = strings.TrimSuffix(queryAddNewVerses, ",\n")

	tx, err := mr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("mr.db.BeginTx: %v", err)
	}
	defer tx.Rollback()

	var id uuid.UUID

	res := tx.QueryRowContext(ctx, queryUpdateSong, upd.GroupName, upd.SongName, song.GroupName, song.SongName)
	if err := res.Scan(&id); err != nil {
		// проверка на уникальность
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("tx.QueryRowContext: %v", err)
	}

	if _, err := tx.ExecContext(ctx, queryUpdateDetail, upd.ReleaseDate, upd.Link, id); err != nil {
		return fmt.Errorf("tx.QueryRowContext: %v", err)
	}

	if _, err := tx.ExecContext(ctx, queryDeleteOldVerses, id); err != nil {
		return fmt.Errorf("tx.QueryRowContext %v", err)
	}

	// опять же, собираем значения для prepared statements
	var vals []interface{}
	for i, verse := range upd.Verses {
		vals = append(vals, id, i+1, verse)
	}

	if _, err := tx.ExecContext(ctx, queryAddNewVerses, vals...); err != nil {
		return fmt.Errorf("tx.QueryRowContext %v", err)
	}

	return tx.Commit()
}

func (mr *MusicRepository) Delete(ctx context.Context, song models.Song) error {
	query :=
		`
	DELETE FROM music_schema.songs AS s
	WHERE s.group_name = $1 AND s.song_name = $2;
	`

	res, err := mr.db.ExecContext(ctx, query, song.GroupName, song.SongName)
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
