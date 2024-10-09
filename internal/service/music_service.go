package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"music_catalog/internal/models"
	"music_catalog/internal/repository/external_api"
	"music_catalog/internal/repository/pg_repo"
)

// musicService is the implementation of the service layer
type musicService struct {
	repo      pg_repo.SongRepository
	apiClient external_api.APIClient
}

// NewMusicService creates a new instance of the MusicService
func NewMusicService(repo pg_repo.SongRepository, apiClient external_api.APIClient) *musicService {
	return &musicService{
		repo:      repo,
		apiClient: apiClient,
	}
}

// AddSong adds a new song to the library and fetches additional song details
func (s *musicService) AddSong(ctx context.Context, group string, title string) error {
	log.Printf("[INFO] Adding song: %s - %s", group, title)

	// Fetch song details from external API
	songDetail, err := s.apiClient.FetchSongDetails(group, title)
	if err != nil {
		log.Printf("[ERROR] Error fetching song details from external API: %v", err)
		return fmt.Errorf("error fetching song details: %w", err)
	}

	log.Printf("[INFO] Fetched song details: %+v", songDetail)

	song, err := s.repo.GetSong(ctx, group, title)
	if err != nil {
		log.Printf("[ERROR] Ошибка при получении песни: %v", err)
		return err
	}

	if song.Title != "" {
		// Песня найдена
		log.Printf("[INFO] Песня уже cуществует: %s - %s", song.Group, song.Title)
		return nil
	}

	// Parse release date
	releaseDate, err := ParseDate(songDetail.ReleaseDate)
	if err != nil {
		log.Printf("[ERROR] Error parsing release date: %v", err)
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
		log.Printf("[ERROR] Error saving song: %v", err)
		return fmt.Errorf("error saving song: %w", err)
	}

	return nil
}

// GetSongs retrieves songs with optional filtering and pagination
func (s *musicService) GetSongs(ctx context.Context, filters models.SongFilters, pagination models.Pagination) ([]models.Song, error) {
	log.Println("[INFO] Fetching songs from library")
	return s.repo.GetSongs(ctx, filters, pagination)
}

// GetSongText retrieves song text with pagination (verse by verse)
func (s *musicService) GetSongText(ctx context.Context, songID int, page int) (string, error) {
	log.Printf("[INFO] Fetching song text for song ID: %d, page: %d", songID, page)
	return s.repo.GetSongText(ctx, songID, page)
}

// UpdateSong updates the details of an existing song
func (s *musicService) UpdateSong(ctx context.Context, song models.Song) error {
	releaseDate, err := ParseDate(song.ReleaseDate)
	if err != nil {
		log.Printf("[ERROR] Error parsing release date: %v", err)
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
