package api

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"strconv"

	catalog_errors "music_catalog/internal/errors"
	"music_catalog/internal/logger"
	"music_catalog/internal/models"

	"github.com/go-chi/chi/v5"
)

// MusicService интерфейс для работы с музыкой
type MusicService interface {
	GetSongs(ctx context.Context, filters models.SongFilters, pagination models.Pagination) ([]models.Song, error)
	AddSong(ctx context.Context, group string, title string) error
	UpdateSong(ctx context.Context, song models.Song) error
	DeleteSong(ctx context.Context, id int) error
	GetSongText(ctx context.Context, songID int, page int) (string, error)
}

// SongHandler handles HTTP requests for songs
type SongHandler struct {
	musicService MusicService
	logger       logger.Logger
}

// NewSongHandler creates a new SongHandler with the provided music service
func NewSongHandler(musicService MusicService, logger logger.Logger) *SongHandler {
	return &SongHandler{
		musicService: musicService,
		logger:       logger,
	}
}

// AddSong adds a new song to the catalog
// @Summary Add a new song
// @Description Adds a new song to the catalog and fetches additional details from an external API
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param song body AddSongRequest true "Song request"
// @Success 201 {string} string "Song added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Error adding the song"
// @Router /songs [post]
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {

	var requestBody AddSongRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		h.logger.Error("Error decoding JSON:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Request to add song: ", requestBody.Group, requestBody.Title)

	if requestBody.Group == "" || requestBody.Title == "" {
		h.logger.Error("Group or song title is missing")
		http.Error(w, "Group or song title is missing", http.StatusBadRequest)
		return
	}

	err = h.musicService.AddSong(r.Context(), requestBody.Group, requestBody.Title)
	if err != nil {
		h.logger.Error("Error adding song:", err)
		http.Error(w, "Error adding the song", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Song added successfully: ", requestBody.Group, requestBody.Title)
	w.WriteHeader(http.StatusCreated)
}

// GetSongs fetches the list of songs with filtering and pagination
// @Summary Get list of songs
// @Description Fetches a list of songs with filtering by all fields and pagination
// @Tags Songs
// @Produce  json
// @Param group query string false "Group name"
// @Param title query string false "Song title"
// @Param release_date query string false "Release date"
// @Param limit query int false "Number of items per page"
// @Param offset query int false "Pagination offset"
// @Success 200 {array} models.Song "List of songs"
// @Failure 500 {string} string "Error retrieving the data"
// @Router /songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры фильтров из запроса
	group := r.URL.Query().Get("group")
	title := r.URL.Query().Get("title")
	releaseDate := r.URL.Query().Get("release_date")

	// Создаем структуру фильтров
	filters := models.SongFilters{
		Group:       group,
		Title:       title,
		ReleaseDate: releaseDate,
	}

	h.logger.Debug("Request to get songs: ", filters.Group, filters.Title, filters.ReleaseDate)

	// Получаем параметры пагинации
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10 // значение по умолчанию
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0 // значение по умолчанию
	}

	pagination := models.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	// Вызов сервиса для получения песен
	songs, err := h.musicService.GetSongs(r.Context(), filters, pagination)
	if err != nil {
		h.logger.Error("Error getting songs:", err)
		http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
		return
	}

	// Формируем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// GetSongText fetches the song lyrics with pagination over verses
// @Summary Get song lyrics with pagination
// @Description Fetches the song lyrics with pagination over verses
// @Tags Songs
// @Produce  json
// @Param id path int true "Song ID"
// @Param page query int false "Page number (default: 0 - full text)"
// @Success 200 {string} string "Paginated verses"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Error retrieving lyrics"
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongText(w http.ResponseWriter, r *http.Request) {

	// Получаем ID песни
	songID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Получаем параметры пагинации из запроса
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)

	if err != nil {
		page = 0 // По умолчанию - полный текст
	}

	h.logger.Debug("Request to get songs", songID, page)

	// Получам полный текст песни
	text, err := h.musicService.GetSongText(r.Context(), songID, page) // 0 обозначает полный текст песни
	if err != nil {
		h.logger.Error("Error getting song text:", err)
		if err == catalog_errors.ErrSongNotFound {
			http.Error(w, "Song not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving lyrics", http.StatusInternalServerError)
		}
		return
	}

	// Формируем JSON ответ
	w.Header().Set("Content-Type", "application/json")

	// Отправляем JSON-ответ с куплетами
	json.NewEncoder(w).Encode(text)

}

// DeleteSong removes a song by ID
// @Summary Delete a song
// @Description Deletes a song by its ID
// @Tags Songs
// @Param id path int true "Song ID"
// @Success 204 {string} string "Song deleted successfully"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Error deleting the song"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	songID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Request to delete song", songID)

	err = h.musicService.DeleteSong(r.Context(), songID)
	if err != nil {
		h.logger.Error("Error deleting song:", err)
		if err == catalog_errors.ErrSongNotFound {
			http.Error(w, "Song not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting the song", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("Song deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Update a song by its ID
// @Description Update an existing song with the provided data.
// @Tags songs
// @Accept  json
// @Produce plain
// @Param id path int true "Song ID"
// @Param song body UpdateSongRequest true "Song data"
// @Success 200 {string} string "Song updated successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Internal server error"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	// Получение ID песни из URL параметров
	songID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		h.logger.Error("Invalid song ID:", err)
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Чтение и декодирование данных из тела запроса
	var updatedSong models.Song
	if err := json.NewDecoder(r.Body).Decode(&updatedSong); err != nil {
		h.logger.Error("Invalid request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updatedSong.Group == "" || updatedSong.Title == "" || updatedSong.Text == "" ||
		updatedSong.Link == "" || updatedSong.ReleaseDate == "" {
		h.logger.Error(fmt.Sprintf("some of required fields is missing: %+v", updatedSong))
		http.Error(w, "some of required fields is missing", http.StatusBadRequest)
		return
	}

	// Добавляем songID в объект песни, чтобы обновить правильную запись
	updatedSong.ID = songID

	h.logger.Debug("Request to update song", updatedSong.ID)

	// Вызов сервиса для обновления песни
	err = h.musicService.UpdateSong(r.Context(), updatedSong)
	if err != nil {
		h.logger.Error("Failed to update song:", err)
		if err == catalog_errors.ErrSongNotFound {
			http.Error(w, "Song not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update song", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("Song updated successfully")
	// Ответ на успешное обновление
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Song updated successfully"))
}

// AddSongRequest данные из запроса для добавления песни в каталог
type AddSongRequest struct {
	Group string `json:"group" example:"Muse"`                   // Artist group
	Title string `json:"song" example:"Supermassive Black Hole"` // Song title
}

// UpdateSongRequest - модель для обновления песни
type UpdateSongRequest struct {
	Group       string `json:"group" example:"Muse"`
	Title       string `json:"title" example:"Supermassive Black Hole"`
	Text        string `json:"text" example:"Ooh baby, don't you know I suffer..."`
	Link        string `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
	ReleaseDate string `json:"release_date" example:"16.07.2006"`
}
