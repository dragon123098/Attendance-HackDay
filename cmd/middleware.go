package main

import "net/http"

func requireRole(role string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. get session from cookie
        // 2. look up user from session
        // 3. if user.Role != role, redirect or 403
        // 4. otherwise call next.ServeHTTP(w, r)
    })
}