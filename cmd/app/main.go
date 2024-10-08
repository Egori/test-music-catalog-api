package main

import (
	"fmt"
	"net/http"

	"music_catalog/config"
	"music_catalog/internal/api"
	"music_catalog/internal/db"
	"music_catalog/internal/logger"
	"music_catalog/internal/repository/external_api"
	"music_catalog/internal/repository/pg_repo"
	"music_catalog/internal/service"

	_ "github.com/lib/pq"

	_ "music_catalog/docs" // Импорт сгенерированных документов
)

func main() {

	logger := logger.NewLogger("debug")

	// Загрузка конфигурации
	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("ошибка при загрузке конфигурации: %v", err)
	}
	logger.Debug(fmt.Sprintf("Loaded config: %+v", config))

	// Подключение к базе данных
	dbConnection, err := db.NewDB(config)
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных: %v", err)
	}
	defer dbConnection.Close()

	// migrations
	err = db.RunMigrations(dbConnection)
	if err != nil {
		logger.Fatal("Ошибка при выполнении миграции: %v", err)
	}

	// Подключаем репозиторий и хендлеры
	repository := pg_repo.NewPostgresSongRepository(dbConnection)
	musicService := service.NewMusicService(repository, external_api.NewExternalAPIClient(config), logger)

	// Инициализация хендлеров
	songHandler := api.NewSongHandler(musicService, logger)

	// Выбираем REST API реализацию
	songAPI := api.NewRestSongAPI(songHandler)

	// Запускаем сервер
	logger.Info(fmt.Sprintf("Starting server on port %s...", config.ServerPort))
	logger.Info(fmt.Sprintf("Swagger UI available at http://localhost:%s/swagger/index.html", config.ServerPort))
	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), songAPI.RegisterRoutes())

}
