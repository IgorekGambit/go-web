package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"

	"go-web/internal/models"
	"go-web/internal/user"
)

type ctxKey int

const userCtxKey ctxKey = iota + 1

const sessionCookieName = "sid"

// UserFromContext возвращает пользователя, установленного миддлварой User.
func UserFromContext(ctx context.Context) (models.User, bool) {
	u, ok := ctx.Value(userCtxKey).(models.User)
	return u, ok
}

func withUser(ctx context.Context, u models.User) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}

// User — cookie sid + EnsureUser в БД; новая сессия получает Set-Cookie.
func User(svc *user.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if svc == nil {
				next.ServeHTTP(w, r)
				return
			}

			c, err := r.Cookie(sessionCookieName)
			sessionID := ""
			if err == nil && c != nil {
				sessionID = c.Value
			}

			created := false
			if sessionID == "" {
				sessionID, err = newSessionID()
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				created = true
			}

			u, err := svc.EnsureUser(r.Context(), sessionID)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if created {
				http.SetCookie(w, sessionCookie(sessionID, r))
			}

			next.ServeHTTP(w, r.WithContext(withUser(r.Context(), u)))
		})
	}
}

func newSessionID() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func sessionCookie(value string, r *http.Request) *http.Cookie {
	secure := r.TLS != nil || os.Getenv("COOKIE_SECURE") == "1"
	return &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	}
}
