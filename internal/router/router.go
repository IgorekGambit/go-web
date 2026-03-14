package router

import (
	"net/http"

	"go-web/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	routes(r)
	return r
}

func routes(r chi.Router) {
	r.Get("/", handlers.Home())
}
