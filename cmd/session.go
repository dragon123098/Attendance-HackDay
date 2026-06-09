package main

import (
	"errors"
	"net/http"
)

func getSessionUser(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", errors.New("no session cookie found")
	}

	username := cookie.Value
	if _, exists := app.Users[username]; !exists {
		return "", errors.New("invalid session user")
	}

	return username, nil
}

func setSessionUser(w http.ResponseWriter, username string) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    username,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func clearSessionUser(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, &cookie)
}