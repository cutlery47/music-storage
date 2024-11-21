package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/cutlery47/music-storage/internal/models"
	"github.com/cutlery47/music-storage/internal/repository"
)

type Service interface {
	Create(ctx context.Context, song models.SongWithDetailPlain) error
	GetSongs(ctx context.Context, limit, offset int, filter models.Filter) ([]models.SongWithDetail, error)
	GetText(ctx context.Context, limit, offset int, song models.Song) (string, error)
	GetDetail(ctx context.Context, song models.Song) (models.SongDetail, error)
	Delete(ctx context.Context, song models.Song) error
	Update(ctx context.Context, song models.Song, upd models.SongWithDetail) error
}

type MusicService struct {
	repo repository.Repository
}

func NewMusicService(repo repository.Repository) *MusicService {
	return &MusicService{
		repo: repo,
	}
}

func (ms *MusicService) Create(ctx context.Context, song models.SongWithDetailPlain) error {
	songSplit := song.Split()
	if err := ms.repo.Create(ctx, songSplit); err != nil {
		return fmt.Errorf("ms.repo.Create: %v", err)
	}
	return nil
}

func (ms *MusicService) GetSongs(ctx context.Context, limit, offset int, filter models.Filter) ([]models.SongWithDetail, error) {
	songs, err := ms.repo.Read(ctx, limit, offset, filter)
	if err != nil {
		return nil, fmt.Errorf("ms.repo.Read: %v", err)
	}
	return songs, nil
}

func (ms *MusicService) GetText(ctx context.Context, limit, offset int, song models.Song) (string, error) {
	verses, err := ms.repo.ReadText(ctx, limit, offset, song)
	if err != nil {
		return "", fmt.Errorf("ms.repo.ReadText: %v", err)
	}

	var res string
	for _, verse := range verses {
		res += fmt.Sprintf("%v\n", verse)
	}
	res = strings.TrimSuffix(res, "\n")

	return res, nil
}

func (ms *MusicService) GetDetail(ctx context.Context, song models.Song) (models.SongDetail, error) {
	detail, err := ms.repo.ReadDetail(ctx, song)
	if err != nil {
		return models.SongDetail{}, fmt.Errorf("ms.repo.ReadDetail: %v", err)
	}
	return detail, nil
}

func (ms *MusicService) Delete(ctx context.Context, song models.Song) error {
	if err := ms.repo.Delete(ctx, song); err != nil {
		return fmt.Errorf("ms.repo.Delete: %v", err)
	}
	return nil
}

func (ms *MusicService) Update(ctx context.Context, song models.Song, upd models.SongWithDetailPlain) error {
	updSplit := upd.Split()
	if err := ms.repo.Update(ctx, song, updSplit); err != nil {
		return fmt.Errorf("ms.repo.Update: %v", err)
	}
	return nil
}
