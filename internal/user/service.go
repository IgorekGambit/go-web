package user

import (
	"context"
	"database/sql"
	"errors"

	"go-web/internal/models"
)

const ensureUserQuery = `
INSERT INTO users (session_id)
VALUES ($1)
ON CONFLICT (session_id) DO UPDATE
SET session_id = EXCLUDED.session_id
RETURNING id, session_id, created_at
`

// Service — создание и выборка пользователя по session_id.
type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	if db == nil {
		return nil
	}
	return &Service{db: db}
}

func (s *Service) EnsureUser(ctx context.Context, sessionID string) (models.User, error) {
	if s == nil || s.db == nil {
		return models.User{}, errors.New("user service: not configured")
	}
	var user models.User
	if err := s.db.QueryRowContext(ctx, ensureUserQuery, sessionID).Scan(
		&user.ID,
		&user.SessionID,
		&user.CreatedAt,
	); err != nil {
		return models.User{}, err
	}
	return user, nil
}
