package handlers

import (
	"net/http"

	"go-web/internal/db"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var version string
		if c := db.DB(); c != nil {
			_ = c.QueryRow("SELECT version()").Scan(&version)
		}
		if version == "" {
			version = "DB not ready"
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("go-web\n\n" + version))
	}
}
