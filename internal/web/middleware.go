package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/PeterGrunig/Attendance-HackDay/internal/domain"
	datastore "github.com/PeterGrunig/Attendance-HackDay/internal/store"
)

type authenticatedUserContextKey struct{}

// RequireRole allows the request when the authenticated user has any permitted role.
func RequireRole(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := loadAuthenticatedUser(w, r)
		if !ok {
			return
		}
		for _, role := range roles {
			if user.Role == role {
				next.ServeHTTP(w, withAuthenticatedUser(r, user))
				return
			}
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := loadAuthenticatedUser(w, r)
		if !ok {
			return
		}
		next.ServeHTTP(w, withAuthenticatedUser(r, user))
	})
}

func loadAuthenticatedUser(w http.ResponseWriter, r *http.Request) (domain.User, bool) {
	userID, err := getSessionUserID(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return domain.User{}, false
	}
	if authStore == nil {
		http.Error(w, "auth store is not configured", http.StatusInternalServerError)
		return domain.User{}, false
	}
	user, err := authStore.FindUserByID(r.Context(), userID)
	if errors.Is(err, datastore.ErrUserNotFound) {
		clearSessionUser(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
		return domain.User{}, false
	}
	if err != nil {
		http.Error(w, "could not load authenticated user", http.StatusInternalServerError)
		return domain.User{}, false
	}
	return user, true
}

func withAuthenticatedUser(r *http.Request, user domain.User) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey{}, user)
	return r.WithContext(ctx)
}

func authenticatedUser(r *http.Request) (domain.User, bool) {
	user, ok := r.Context().Value(authenticatedUserContextKey{}).(domain.User)
	return user, ok
}
