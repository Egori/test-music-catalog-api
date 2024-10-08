package api

import (
	"net/http"
)

type SongAPI interface {
	RegisterRoutes() http.Handler
}
