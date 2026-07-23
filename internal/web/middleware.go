package web

import (
	"context"
	"errors"
	"log"
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
		log.Printf("authorization denied: method=%s path=%q role=%q permitted_roles=%v", r.Method, r.URL.Path, user.Role, roles)
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
		log.Printf("authentication failed: method=%s path=%q error=%v", r.Method, r.URL.Path, err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return domain.User{}, false
	}
	if authStore == nil {
		log.Printf("authentication failed: method=%s path=%q user_id=%q error=auth store is not configured", r.Method, r.URL.Path, userID)
		http.Error(w, "auth store is not configured", http.StatusInternalServerError)
		return domain.User{}, false
	}
	user, err := authStore.FindUserByID(r.Context(), userID)
	if errors.Is(err, datastore.ErrUserNotFound) {
		log.Printf("authentication failed: method=%s path=%q user_id=%q error=%v", r.Method, r.URL.Path, userID, err)
		clearSessionUser(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
		return domain.User{}, false
	}
	if err != nil {
		log.Printf("authentication failed: method=%s path=%q user_id=%q error=%v", r.Method, r.URL.Path, userID, err)
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
