package handlers

import (
	"log"
	"net/http"
	"time"

	"go-web/internal/db"
	"go-web/resources/views"
)

type homeData struct {
	Title      string
	DBVersion  string
	Hint       string
	RenderedAt string
}

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var version string
		if c := db.DB(); c != nil {
			_ = c.QueryRow("SELECT version()").Scan(&version)
		}
		if version == "" {
			version = "DB not ready"
		}

		data := homeData{
			Title:      "Main Page",
			DBVersion:  version,
			Hint:       "Поля Title, Hint, RenderedAt передаются из Go в шаблоны: layout (base) и страница (content).",
			RenderedAt: time.Now().Format(time.RFC3339),
		}

		if err := views.RenderHTML(w, "home", data); err != nil {
			log.Printf("render home: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
