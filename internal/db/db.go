package db

import (
	"database/sql"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	conn *sql.DB
	once sync.Once
)

// Open открывает подключение к БД и сохраняет его в пакете.
func Open(dsn string) error {
	var err error
	once.Do(func() {
		conn, err = sql.Open("pgx", dsn)
	})
	return err
}

// Close закрывает подключение.
func Close() error {
	if conn == nil {
		return nil
	}
	return conn.Close()
}

// DB возвращает текущее подключение. Вызывать после Open().
func DB() *sql.DB {
	return conn
}
