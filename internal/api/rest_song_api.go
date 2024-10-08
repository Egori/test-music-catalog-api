package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type RestSongAPI struct {
	songHandler *SongHandler
}

func NewRestSongAPI(songHandler *SongHandler) *RestSongAPI {
	return &RestSongAPI{songHandler: songHandler}
}

func (api *RestSongAPI) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/songs", api.songHandler.GetSongs)
	r.Get("/songs/{id}/text", api.songHandler.GetSongText)
	r.Post("/songs", api.songHandler.AddSong)
	r.Put("/songs/{id}", api.songHandler.UpdateSong)
	r.Delete("/songs/{id}", api.songHandler.DeleteSong)
	// Маршрут для Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	return r
}
