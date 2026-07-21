package web

import (
	"errors"
	"log"
	"net/http"
	"strings"

	datastore "github.com/PeterGrunig/Attendance-HackDay/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// loginHandler renders the login page and traces each authentication stage without logging secrets.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	//log.Printf("login request: method=%s path=%q", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		//log.Printf("login page requested: path=%q", r.URL.Path)
		renderUnAuth(w, "login.html", PageData{Title: "Login"})
		return

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Printf("login failed: path=%q stage=parse_form error=%v", r.URL.Path, err)
			http.Error(w, "invalid form submission", http.StatusBadRequest)
			return
		}

		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")
		//log.Printf("login submitted: path=%q email=%q password_provided=%t", r.URL.Path, email, password != "")

		if email == "" || password == "" {
			//log.Printf("login rejected: path=%q email=%q stage=validate email_provided=%t password_provided=%t", r.URL.Path, email, email != "", password != "")
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Email and password are required.",
			})
			return
		}

		if authStore == nil {
			log.Printf("login failed: path=%q email=%q stage=store_lookup error=auth store is not configured", r.URL.Path, email)
			http.Error(w, "auth store is not configured", http.StatusInternalServerError)
			return
		}

		log.Printf("login user lookup started: path=%q email=%q", r.URL.Path, email)
		user, err := authStore.FindUserByEmail(r.Context(), email)
		if errors.Is(err, datastore.ErrUserNotFound) {
			log.Printf("login rejected: path=%q email=%q stage=store_lookup reason=user_not_found", r.URL.Path, email)
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Invalid email or password.",
			})
			return
		}
		if err != nil {
			log.Printf("login failed: path=%q email=%q stage=store_lookup error=%v", r.URL.Path, email, err)
			http.Error(w, "could not authenticate user", http.StatusInternalServerError)
			return
		}
		log.Printf("login user found: path=%q email=%q user_id=%q role=%q", r.URL.Path, email, user.UserID, user.Role)

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			log.Printf("login rejected: path=%q email=%q user_id=%q stage=password_check error=%v", r.URL.Path, email, user.UserID, err)
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Invalid email or password.",
			})
			return
		}
		log.Printf("login password verified: path=%q email=%q user_id=%q", r.URL.Path, email, user.UserID)

		if err := createSession(w, user.UserID); err != nil {
			log.Printf("login failed: path=%q email=%q user_id=%q stage=create_session error=%v", r.URL.Path, email, user.UserID, err)
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}
		log.Printf("login session created: path=%q user_id=%q role=%q", r.URL.Path, user.UserID, user.Role)

		switch user.Role {
		case "student":
			log.Printf("login succeeded: user_id=%q role=%q redirect=%q", user.UserID, user.Role, "/studentDashboard")
			http.Redirect(w, r, "/studentDashboard", http.StatusSeeOther)
		case "teacher":
			log.Printf("login succeeded: user_id=%q role=%q redirect=%q", user.UserID, user.Role, "/teacherDashboard")
			http.Redirect(w, r, "/teacherDashboard", http.StatusSeeOther)
		case "admin":
			log.Printf("login succeeded: user_id=%q role=%q redirect=%q", user.UserID, user.Role, "/adminDashboard")
			http.Redirect(w, r, "/adminDashboard", http.StatusSeeOther)
		default:
			log.Printf("login rejected: user_id=%q role=%q stage=role_redirect reason=invalid_role", user.UserID, user.Role)
			http.Error(w, "invalid user role", http.StatusForbidden)
		}

		return

	default:
		log.Printf("login rejected: method=%s path=%q reason=method_not_allowed", r.Method, r.URL.Path)
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
