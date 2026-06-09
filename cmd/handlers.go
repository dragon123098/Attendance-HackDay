package main

import (
	"html/template"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

// PageData holds transient page-only values.
type PageData struct {
	Title       string
	Username    string
	AvatarImage string
	Error       string
}

func loginView(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		data := PageData{
			Title: "Login",
		}
		renderUnAuth(w, "login.html", data)
		return

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form submission", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Email and password are required.",
			})
			return
		}

		var user *User

		for _, u := range app.Users {
			if u.Email == email {
				user = u
				break
			}
		}

		if user == nil {
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Invalid email or password.",
			})
			return
		}

		err := bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(password),
		)

		if err != nil {
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Invalid email or password.",
			})
			return
		}

		setSessionUser(w, user.UserID)

		switch user.Role {
		case "student":
			http.Redirect(w, r, "/studentDashboard", http.StatusSeeOther)

		case "teacher":
			http.Redirect(w, r, "/teacherDashboard", http.StatusSeeOther)

		case "admin":
			http.Redirect(w, r, "/adminDashboard", http.StatusSeeOther)

		default:
			clearSessionUser(w)
			http.Error(w, "invalid user role", http.StatusForbidden)
		}

		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func logoutView(w http.ResponseWriter, r *http.Request) {
	clearSessionUser(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func studentView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user := app.Users[username]
	data := PageData{
		Title:       "Student Dashboard",
		Username:    user.Name,
		AvatarImage: "/static/images/geraldIcon3.png",
	}
	render(w, "studentDash.html", data)
}

func shopView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user := app.Users[username]
	data := PageData{
		Title:       "Shop",
		Username:    user.Name,
		AvatarImage: "/static/images/geraldIcon3.png",
	}
	render(w, "shopView.html", data)
}

func avatarView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user := app.Users[username]
	data := PageData{
		Title:       "Avatar",
		Username:    user.Name,
		AvatarImage: "/static/images/geraldIcon3.png",
	}
	render(w, "avatarView.html", data)
}

func coinView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user := app.Users[username]
	data := PageData{
		Title:    "Coins",
		Username: user.Name,
	}
	render(w, "coinView.html", data)
}

func teacherView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user := app.Users[username]
	data := PageData{
		Title:    "Teacher Dashboard",
		Username: user.Name,
	}
	render(w, "teacherDash.html", data)
}

func teacherEditView(w http.ResponseWriter, r *http.Request) {
	render(w, "teacherEdit.html", nil)
}

func adminView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user := app.Users[username]
	data := PageData{
		Title:    "Admin Dashboard",
		Username: user.Name,
	}
	render(w, "adminDash.html", data)
}

func adminEditView(w http.ResponseWriter, r *http.Request) {
	render(w, "adminEdit.html", nil)
}

func render(w http.ResponseWriter, page string, data any) {
	tmpl, err := loadTemplates(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadTemplates(page string) (*template.Template, error) {
	return template.ParseFiles(
		filepath.Join("templates", "AuthBase.html"),
		filepath.Join("templates", "partials", "topbar.html"),
		filepath.Join("templates", "partials", "navbar.html"),
		filepath.Join("templates", "partials", "footer.html"),
		filepath.Join("templates", page),
	)
}

func loadUnAuthTemplates(page string) (*template.Template, error) {
	return template.ParseFiles(
		filepath.Join("templates", "UnAuthBase.html"),
		filepath.Join("templates", page),
	)
}

func renderUnAuth(w http.ResponseWriter, page string, data any) {
	tmpl, err := loadUnAuthTemplates(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
