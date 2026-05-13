package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go-web/internal/db"
	"go-web/internal/migrate"
	"go-web/internal/router"
)

func main() {
	dsn := getEnv("DB_DSN", "postgres://user:pass@localhost:5432/app")
	port := getEnv("PORT", "8080")

	if err := db.Open(dsn); err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	if c := db.DB(); c != nil {
		if err := c.Ping(); err != nil {
			log.Printf("db ping (возможно БД ещё не готова): %v", err)
		}
		if err := migrate.Up(context.Background(), c); err != nil {
			log.Fatalf("migrate: %v", err)
		}
	}

	r := router.New(db.DB())

	log.Printf("Listening on :%s", port)
	http.ListenAndServe(":"+port, r)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
