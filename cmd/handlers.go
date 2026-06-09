package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func loginView(w http.ResponseWriter, r *http.Request) {
	render(w, "home.html", nil)
}

func loginSubmitView(w http.ResponseWriter, r *http.Request) {
}

func logoutView(w http.ResponseWriter, r *http.Request) {
}

func studentView(w http.ResponseWriter, r *http.Request) {
}

func shopView(w http.ResponseWriter, r *http.Request) {
}

func avatarView(w http.ResponseWriter, r *http.Request) {
}

func coinView(w http.ResponseWriter, r *http.Request) {
}

func teacherView(w http.ResponseWriter, r *http.Request) {
}

func teacherEditView(w http.ResponseWriter, r *http.Request) {
}

func adminView(w http.ResponseWriter, r *http.Request) {
}

func adminEditView(w http.ResponseWriter, r *http.Request) {
}

// render is a helper function to render templates with a base layout
func render(w http.ResponseWriter, page string, data any) {
	tmpl, err := loadTemplates(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


//This will load the templates needed. Just pass in what page you want and it will render it will all the correct stuff
func loadTemplates(page string) (*template.Template, error) {
	return template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", "partials", "topbar.html"),
		filepath.Join("templates", "partials", "navbar.html"),
		filepath.Join("templates", "partials", "footer.html"),
		filepath.Join("templates", page),
	)
}