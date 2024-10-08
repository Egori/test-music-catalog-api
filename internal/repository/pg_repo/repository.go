package pg_repo

import (
	"context"
	"music_catalog/internal/models"
)

// SongRepository — интерфейс для работы с репозиторием песен.
type SongRepository interface {
	GetSongText(ctx context.Context, songID int, page int) (string, error)                                         // Получить текст песни по ID
	GetSongs(ctx context.Context, filters models.SongFilters, pagination models.Pagination) ([]models.Song, error) // Получить список песен с фильтрацией и пагинацией
	AddSong(ctx context.Context, song models.Song) (int, error)                                                    // Добавить новую песню
	GetSongByID(ctx context.Context, id int) (models.Song, error)                                                  // Получить песню по ID
	GetSong(ctx context.Context, group string, title string) (models.Song, error)                                  // Получить песню по group_name и title                                          // Получить песню по ID
	UpdateSong(ctx context.Context, song models.Song) error                                                        // Обновить песню
	DeleteSong(ctx context.Context, id int) error                                                                  // Удалить песню
}
