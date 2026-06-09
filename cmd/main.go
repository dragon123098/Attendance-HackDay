package main

import (
	// native Go packages
	"log"
	"net/http"
	// internal packages
	// 3rd party packages
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	//all roles
	mux.HandleFunc("GET /login", loginView)              	// Returns full page, login view (initial page)
	mux.HandleFunc("POST /login", loginSubmitView)       	// Returns HTML fragment (HTMX), login submission view
	mux.HandleFunc("POST /logout", logoutView)       		// Returns HTML fragment (HTMX), logout view

	//student role
	mux.Handle("GET /studentDashboard", requireRole("student", http.HandlerFunc(studentView)))			// Returns full page, student dashboard view
	mux.HandleFunc("GET /shop", shopView)        		 	// Returns full page, shop view
	mux.HandleFunc("GET /avatar", avatarView)        	 	// Returns full page, avatar view
	mux.HandleFunc("POST /studentDashboard/coin", coinView) // Returns HTML fragment (HTMX), coin popup after marking attendance
	//teacher role
	mux.Handle("GET /teacherDashboard", requireRole("teacher", http.HandlerFunc(teacherView)))			// Returns full page, teacher dashboard view
	mux.HandleFunc("POST /teacherDashboard/edit", teacherEditView) // Returns HTML fragment (HTMX), edit view
	//admin role
	mux.Handle("GET /adminDashboard", requireRole("admin", http.HandlerFunc(adminView)))				// Returns full page, admin dashboard view
	mux.HandleFunc("POST /adminDashboard/edit", adminEditView)   // Returns HTML fragment (HTMX), edit view

	log.Print("starting server on http://localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
