package main

import "net/http"

func requireRole(role string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. look up user from session
		username, err := getSessionUser(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
        // 2. if user.Role != role, redirect or 403
		user := app.Users[username]
        if user.Role != role {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        // 3. otherwise call next.ServeHTTP(w, r)
		next.ServeHTTP(w, r)
    })
}

func requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. look up user from session
		_, err := getSessionUser(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// 2. if no user, redirect to login
		// 3. otherwise call next.ServeHTTP(w, r)
		next.ServeHTTP(w, r)
	})
}