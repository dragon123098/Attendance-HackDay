package web

import (
	"errors"
	"net/http"
	"strings"

	datastore "github.com/dragon123098/Attendance-HackDay.git/internal/store"
	"golang.org/x/crypto/bcrypt"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderUnAuth(w, "login.html", PageData{Title: "Login"})
		return

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form submission", http.StatusBadRequest)
			return
		}

		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		if email == "" || password == "" {
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Email and password are required.",
			})
			return
		}

		if authStore == nil {
			http.Error(w, "auth store is not configured", http.StatusInternalServerError)
			return
		}

		user, err := authStore.FindUserByEmail(r.Context(), email)
		if errors.Is(err, datastore.ErrUserNotFound) {
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Invalid email or password.",
			})
			return
		}
		if err != nil {
			http.Error(w, "could not authenticate user", http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Invalid email or password.",
			})
			return
		}

		if err := createSession(w, user.UserID); err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		switch user.Role {
		case "student":
			http.Redirect(w, r, "/studentDashboard", http.StatusSeeOther)
		case "teacher":
			http.Redirect(w, r, "/teacherDashboard", http.StatusSeeOther)
		case "admin":
			http.Redirect(w, r, "/adminDashboard", http.StatusSeeOther)
		default:
			http.Error(w, "invalid user role", http.StatusForbidden)
		}

		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func logoutView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clearSessionUser(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
