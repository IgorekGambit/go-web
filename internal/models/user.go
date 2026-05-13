package models

import (
	"time"
)

// User — доменная модель пользователя (сессия без профиля).
type User struct {
	ID        string
	SessionID string
	CreatedAt time.Time
}
