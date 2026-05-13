package router

import (
	"database/sql"
	"net/http"

	"go-web/internal/handlers"
	appmw "go-web/internal/middleware"
	"go-web/internal/user"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func New(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	if db != nil {
		r.Use(appmw.User(user.NewService(db)))
	}
	routes(r)
	return r
}

func routes(r chi.Router) {
	r.Get("/", handlers.Home())
}
