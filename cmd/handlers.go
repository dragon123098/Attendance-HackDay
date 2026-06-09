package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

//idk if this is the right place to this, but I'm putting it here
//Any pagedata we ned to add to the page, we populate this struct with data from the database
//This is basically a non persistable entity

type PageData struct {
	Title       string
	Username    string
	AvatarImage string
}

func loginView(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Login",
	}

	renderUnAuth(w, "login.html", data)
}

func loginSubmitView(w http.ResponseWriter, r *http.Request) {
}

func logoutView(w http.ResponseWriter, r *http.Request) {
}

func studentView(w http.ResponseWriter, r *http.Request) {
}

func shopView(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		AvatarImage: "/static/images/geraldIcon2.png",
	}

	render(w, "home.html", data)
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


//idk if this is the right place for this, but it makes sense to me
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
		filepath.Join("templates", "AuthBase.html"),
		filepath.Join("templates", "partials", "topbar.html"),
		filepath.Join("templates", "partials", "navbar.html"),
		filepath.Join("templates", "partials", "footer.html"),
		filepath.Join("templates", page),
	)
}


//These two functions load pages for unauthenticated users.
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

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}