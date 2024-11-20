package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cutlery47/music-storage/internal/models"
	"github.com/cutlery47/music-storage/internal/repository"
	"github.com/cutlery47/music-storage/internal/service"
)

func Run() error {
	ctx := context.Background()

	url := "postgresql://postgres:postgres@localhost:5433/music?sslmode=disable"

	repo, err := repository.NewMusicRepository(url)
	if err != nil {
		return fmt.Errorf("couldn't connect to the database: %v", err)
	}

	srv := service.NewMusicService(repo)

	date := time.Date(2001, 10, 8, 0, 0, 0, 0, time.UTC)
	group := "Linkin Park"

	filter := models.Filter{
		Group:         &group,
		ReleasedAfter: &date,
	}

	songs, err := srv.GetSongs(ctx, 0, 10, filter)
	if err != nil {
		return fmt.Errorf("srv.GetSongs: %v", err)
	}

	log.Println("songs:", songs)

	song := models.Song{
		GroupName: "Linkin Park",
		SongName:  "In the End",
	}

	info, err := srv.GetDetail(ctx, song)
	if err != nil {
		return fmt.Errorf("srv.GetDetail: %v", err)
	}

	log.Println("info:", info)

	// if err := srv.Delete(ctx, song); err != nil {
	// 	return fmt.Errorf("srv.Delete: %v", err)
	// }

	return nil
}
