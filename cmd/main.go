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
	// auth routes
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("POST /logout", logoutView)
	// student routes
	mux.Handle("GET /studentDashboard", requireRole("student", http.HandlerFunc(studentView)))
	mux.Handle("GET /shop", requireRole("student", http.HandlerFunc(shopView)))
	mux.Handle("GET /avatar", requireRole("student", http.HandlerFunc(avatarView)))
	// teacher routes
	mux.Handle("GET /teacherDashboard", requireRole("teacher", http.HandlerFunc(teacherView)))
	mux.HandleFunc("POST /teacherDashboard/edit", teacherEditView)
	// admin routes
	mux.Handle("GET /adminDashboard", requireRole("admin", http.HandlerFunc(adminView)))
	mux.HandleFunc("POST /adminDashboard/edit", adminEditView)

	log.Print("starting server on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
