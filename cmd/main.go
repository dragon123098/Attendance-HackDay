package main

import (
	"log"
	"net/http"
)

var app AppState


func main() {
	loadData()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	//all roles
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("POST /logout", logoutView)
	//student role
	mux.Handle("GET /studentDashboard", requireRole("student", http.HandlerFunc(studentView)))
	mux.HandleFunc("GET /shop", shopView)
	mux.HandleFunc("GET /avatar", avatarView)
	//teacher role
	mux.Handle("GET /teacherDashboard", requireRole("teacher", http.HandlerFunc(teacherView)))
	mux.HandleFunc("POST /teacherDashboard/edit", teacherEditView)
	//admin role
	mux.Handle("GET /adminDashboard", requireRole("admin", http.HandlerFunc(adminView)))
	mux.HandleFunc("POST /adminDashboard/edit", adminEditView)

	log.Print("starting server on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
