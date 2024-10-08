package models

// Songs - структура для хранения данных о песне
type Song struct {
	ID          int    `json:"id"`
	Group       string `json:"group"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	Link        string `json:"link"`
	ReleaseDate string `json:"release_date"`
}

// SongFilters - структура для хранения фильтров для запросов к базе данных
type SongFilters struct {
	Group       string // фильтр по названию группы
	Title       string // фильтр по названию песни
	ReleaseDate string // фильтр по дате выпуска
}

// Pagination - структура для хранения параметров пагинации
type Pagination struct {
	Limit  int // максимальное количество записей на странице
	Offset int // смещение (номер записи, с которой начинать выборку)
}
