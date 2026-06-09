package main

import (
	"net/http"
)

func studentView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, ok := app.Users[username]
	if !ok {
		clearSessionUser(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := PageData{
		Title:       "Student Dashboard",
		Username:    user.Name,
		AvatarImage: "/static/images/geraldIcon3.png",
		Coins: 100,
	}

	renderStudent(w, "studentDash.html", data)
}

func shopView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, ok := app.Users[username]
	if !ok {
		clearSessionUser(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := PageData {
		Title: "Shop",
		Username:    user.Name,
		AvatarImage:  "/static/images/geraldIcon3.png",
		Coins: 100,
	}

	renderStudent(w, "shopView.html", data)
}


func avatarView(w http.ResponseWriter, r *http.Request) {
	username, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, ok := app.Users[username]
	if !ok {
		clearSessionUser(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := PageData{
		Title:       "Avatar",
		Username:    user.Name,
		AvatarImage: "/static/Images/gerald.png",
		Coins:       42,
	}


	renderStudent(w, "avatarView.html", data)
}