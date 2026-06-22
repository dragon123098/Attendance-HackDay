package web

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"sync"
	"time"
)

type sessionRecord struct {
	UserID    string
	ExpiresAt time.Time
}

var (
	sessionMu    sync.RWMutex
	sessionStore = map[string]sessionRecord{}
)

const sessionDuration = 24 * time.Hour

func createSession(w http.ResponseWriter, userID string) error {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	sessionMu.Lock()
	sessionStore[token] = sessionRecord{
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	sessionMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionDuration.Seconds()),
	})

	return nil
}

func getSessionUser(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", errors.New("no session cookie found")
	}

	sessionMu.RLock()
	sess, ok := sessionStore[cookie.Value]
	sessionMu.RUnlock()

	if !ok {
		return "", errors.New("invalid session")
	}

	if time.Now().After(sess.ExpiresAt) {
		sessionMu.Lock()
		delete(sessionStore, cookie.Value)
		sessionMu.Unlock()
		return "", errors.New("session expired")
	}

	if _, exists := app.Users[sess.UserID]; !exists {
		return "", errors.New("invalid session user")
	}

	return sess.UserID, nil
}

func clearSessionUser(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil {
		sessionMu.Lock()
		delete(sessionStore, cookie.Value)
		sessionMu.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
