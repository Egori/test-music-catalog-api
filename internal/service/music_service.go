package service

import (
	"context"
	"fmt"
	"log"
	"time"

	catalog_errors "music_catalog/internal/errors"
	"music_catalog/internal/logger"
	"music_catalog/internal/models"
	"music_catalog/internal/repository/external_api"
	"music_catalog/internal/repository/pg_repo"
)

// musicService is the implementation of the service layer
type musicService struct {
	repo      pg_repo.SongRepository
	apiClient external_api.APIClient
	logger    logger.Logger
}

// NewMusicService creates a new instance of the MusicService
func NewMusicService(repo pg_repo.SongRepository, apiClient external_api.APIClient, logger logger.Logger) *musicService {
	return &musicService{
		repo:      repo,
		apiClient: apiClient,
		logger:    logger,
	}
}

// AddSong adds a new song to the library and fetches additional song details
func (s *musicService) AddSong(ctx context.Context, group string, title string) error {

	// Fetch song details from external API
	songDetail, err := s.apiClient.FetchSongDetails(group, title)
	if err != nil {
		log.Printf("[ERROR] Error fetching song details from external API: %v", err)
		return fmt.Errorf("error fetching song details: %w", err)
	}

	s.logger.Info("Song found in external API")
	s.logger.Debug(fmt.Sprintf("Song details from external API: %+v", *songDetail))

	song, err := s.repo.GetSong(ctx, group, title)
	if err != nil {
		s.logger.Error("Error getting song from repository: ", err)
		return err
	}

	if song.Title != "" {
		// Песня найдена
		s.logger.Info("Song found in library: ", song.Title)
		return catalog_errors.ErrSongExists
	}

	// Parse release date
	releaseDate, err := ParseDate(songDetail.ReleaseDate)
	if err != nil {
		s.logger.Error("Error parsing release date: ", err)
		return fmt.Errorf("error parsing release date: %w", err)
	}
	songDetail.ReleaseDate = releaseDate.Format("2006-01-02")

	// Create new song instance
	newSong := models.Song{
		Group:       group,
		Title:       title,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	// Save to repository
	if _, err := s.repo.AddSong(ctx, newSong); err != nil {
		s.logger.Error("Error saving song: ", err)
		return fmt.Errorf("error saving song: %w", err)
	}

	return nil
}

// GetSongs retrieves songs with optional filtering and pagination
func (s *musicService) GetSongs(ctx context.Context, filters models.SongFilters, pagination models.Pagination) ([]models.Song, error) {

	if filters.ReleaseDate != "" {
		releaseDate, err := ParseDate(filters.ReleaseDate)
		if err != nil {
			s.logger.Error("Error parsing release date: ", err)
			return []models.Song{}, fmt.Errorf("error parsing release date: %w", err)
		}
		filters.ReleaseDate = releaseDate.Format("2006-01-02")
	}
	return s.repo.GetSongs(ctx, filters, pagination)
}

// GetSongText retrieves song text with pagination (verse by verse)
func (s *musicService) GetSongText(ctx context.Context, songID int, page int) (string, error) {
	return s.repo.GetSongText(ctx, songID, page)
}

// UpdateSong updates the details of an existing song
func (s *musicService) UpdateSong(ctx context.Context, song models.Song) error {
	releaseDate, err := ParseDate(song.ReleaseDate)
	if err != nil {
		s.logger.Error("Error parsing release date: ", err)
		return fmt.Errorf("error parsing release date: %w", err)
	}
	song.ReleaseDate = releaseDate.Format("2006-01-02")
	return s.repo.UpdateSong(ctx, song)
}

// DeleteSong deletes a song from the library
func (s *musicService) DeleteSong(ctx context.Context, songID int) error {
	return s.repo.DeleteSong(ctx, songID)
}

// Список возможных форматов дат
var dateFormats = []string{
	"02.01.2006",      // Формат "DD.MM.YYYY"
	"2006-01-02",      // Формат "YYYY-MM-DD"
	"January 2, 2006", // Формат "January 2, 2006"
	time.RFC3339,      // ISO 8601 формат (например, "2006-01-02T15:04:05Z07:00")
}

// ParseDate пытается разобрать дату в одном из поддерживаемых форматов
func ParseDate(dateStr string) (time.Time, error) {
	var parsedDate time.Time
	var err error

	// Перебор всех возможных форматов
	for _, layout := range dateFormats {
		parsedDate, err = time.Parse(layout, dateStr)
		if err == nil {
			// Если дата успешно разобрана в одном из форматов — возвращаем её
			return parsedDate, nil
		}
	}

	// Если ни один формат не подошел, возвращаем ошибку
	return time.Time{}, fmt.Errorf("не удалось распознать формат даты: %s", dateStr)
}
