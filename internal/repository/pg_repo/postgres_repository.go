package pg_repo

import (
	"context"
	"database/sql"
	"fmt"
	catalog_errors "music_catalog/internal/errors"
	"music_catalog/internal/models"
	"strings"
)

// PostgresMusicRepository — структура для работы с PostgreSQL.
type PostgresMusicRepository struct {
	db *sql.DB
}

// NewPostgresSongRepository — конструктор для PostgresSongRepository.
func NewPostgresSongRepository(db *sql.DB) *PostgresMusicRepository {
	return &PostgresMusicRepository{db: db}
}

// GetSongs — получение списка песен с фильтрацией и пагинацией.
func (r *PostgresMusicRepository) GetSongs(ctx context.Context, filters models.SongFilters, pagination models.Pagination) ([]models.Song, error) {
	query := "SELECT id, group_name, title, release_date, text, link FROM songs WHERE 1=1" // базовый запрос
	args := []interface{}{}
	argCount := 1

	// Фильтрация по группе
	if filters.Group != "" {
		query += fmt.Sprintf(" AND group_name ILIKE $%d", argCount) // динамически добавляем фильтр
		args = append(args, "%"+filters.Group+"%")
		argCount++
	}

	// Фильтрация по названию песни
	if filters.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argCount)
		args = append(args, "%"+filters.Title+"%")
		argCount++
	}

	// Фильтрация по дате выпуска
	if filters.ReleaseDate != "" {
		query += fmt.Sprintf(" AND release_date = $%d", argCount)
		args = append(args, filters.ReleaseDate)
		argCount++
	}

	// Пагинация добавляется в конец запроса
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, pagination.Limit, pagination.Offset)

	// Выполнение запроса
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Парсинг результатов
	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Title, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

// AddSong — добавление новой песни в базу данных.
func (r *PostgresMusicRepository) AddSong(ctx context.Context, song models.Song) (int, error) {
	var id int
	query := `INSERT INTO songs (group_name, title, text, link, release_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, song.Group, song.Title, song.Text, song.Link, song.ReleaseDate).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении песни: %w", err)
	}
	return id, nil
}

// GetSongByID — получение песни по ID.
func (r *PostgresMusicRepository) GetSongByID(ctx context.Context, id int) (models.Song, error) {
	var song models.Song
	query := `SELECT id, group_name, title, text, link, release_date FROM songs WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&song.ID, &song.Group, &song.Title, &song.Text, &song.Link, &song.ReleaseDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Song{}, nil // Вернем пустую песню, если запись не найдена
		}
		return models.Song{}, fmt.Errorf("ошибка при получении песни по ID: %w", err)
	}
	return song, nil
}

// GetSong — получение песни по group_name и title.
func (r *PostgresMusicRepository) GetSong(ctx context.Context, group string, title string) (models.Song, error) {
	var song models.Song
	query := `SELECT id, group_name, title, text, link, release_date FROM songs WHERE group_name = $1 AND title = $2`
	err := r.db.QueryRowContext(ctx, query, group, title).Scan(&song.ID, &song.Group, &song.Title, &song.Text, &song.Link, &song.ReleaseDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Song{}, nil // Вернем пустую песню, если запись не найдена
		}
		return models.Song{}, fmt.Errorf("ошибка при получении песни: %w", err)
	}
	return song, nil
}

// UpdateSong — обновление данных песни.
func (r *PostgresMusicRepository) UpdateSong(ctx context.Context, song models.Song) error {
	query := `UPDATE songs SET group_name = $1, title = $2, text = $3, link = $4, release_date = $5, updated_at = NOW() WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, song.Group, song.Title, song.Text, song.Link, song.ReleaseDate, song.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении песни: %w", err)
	}
	return nil
}

// DeleteSong — удаление песни по ID.
func (r *PostgresMusicRepository) DeleteSong(ctx context.Context, id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении песни: %w", err)
	}
	return nil
}

// GetSongText retrieves song text with pagination (verse by verse)
func (r *PostgresMusicRepository) GetSongText(ctx context.Context, songID int, page int) (string, error) {
	// Query to get the text for the song
	query := `SELECT text FROM songs WHERE id = $1`

	var fullText string

	// Execute query
	err := r.db.QueryRowContext(ctx, query, songID).Scan(&fullText)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", catalog_errors.ErrSongNotFound
		}
		return "", fmt.Errorf("error fetching song text: %w", err)
	}

	if page == 0 {
		return fullText, nil
	}

	// Split the text into verses (assuming each verse is separated by a newline)
	verses := strings.Split(fullText, "\n\n")
	if page < 1 || page > len(verses) {
		return "", catalog_errors.ErrInvalidPage
	}

	// Return the specific verse for the given page
	return verses[page-1], nil
}
