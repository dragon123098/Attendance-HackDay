package web

import "net/http"

func RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, err := getSessionUser(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, ok := app.Users[username]
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if user.Role != role {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := getSessionUser(r); err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}