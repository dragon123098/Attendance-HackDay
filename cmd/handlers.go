package main

import (
	"net/http"
	"log"
)

//idk if this is the right place to this, but I'm putting it here
//Any pagedata we ned to add to the page, we populate this struct with data from the database
//This is basically a non persistable entity

type PageData struct {
	Title       string
	Username    string
	AvatarImage string
	Coins int
	Items []Shop
}

//This is just for testing, we can keep it and populate it form the database, or just put the DB directly to the page
type Shop struct {
	Name string
	Cost int
	Image string
}



func loginView(w http.ResponseWriter, r *http.Request) {
		data := PageData{
		Title: "Login",
	}
	renderUnAuth(w, "login.html", data)
}

func loginPostView(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Printf("Username: %s", username)
	log.Printf("Password: %s", password)

	http.Redirect(w, r, "/shop", http.StatusSeeOther)
}

func logoutView(w http.ResponseWriter, r *http.Request) {
	renderUnAuth(w, "logout.html", nil)
}

func studentView(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		AvatarImage: "/static/images/geraldIcon3.png",
		Coins: 100,
	}
	renderStudent(w, "studentDash.html", data)

}

func shopView(w http.ResponseWriter, r *http.Request) {
	data := PageData {
		Title: "Shop",
		AvatarImage: "/static/images/geraldIcon3.png",
		Coins: 100,
		Items: []Shop{},
	}
	renderStudent(w, "shopView.html", data)
}


func avatarView(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:       "Avatar",
		AvatarImage: "/static/Images/gerald.png",
		Coins:       42,
	}

	renderStudent(w, "avatarView.html", data)
}


func teacherView(w http.ResponseWriter, r *http.Request) {
	renderTeacher(w, "teacherDash.html", nil)
}

func teacherEditView(w http.ResponseWriter, r *http.Request) {
	renderTeacher(w, "teacherEdit.html", nil)
}

func adminView(w http.ResponseWriter, r *http.Request) {
	renderAdmin(w, "adminDash.html", nil)
}

func adminEditView(w http.ResponseWriter, r *http.Request) {
	renderAdmin(w, "adminEdit.html", nil)
}



