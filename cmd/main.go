package main

import (
	"log"
	"net/http"
)

var app AppState


func main() {
	loadData()
	app = AppState{
    Users:         make(map[string]*User),
    Classrooms:    make(map[string]*Classroom),
    ShopItems:     make(map[string]*ShopItem),
    AvatarConfigs: make(map[string]*AvatarConfig),
    Transactions:  []CoinTransaction{},
    Attendance:    []AttendanceRecord{},
    Schedule:      []Schedule{},
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	//all roles
	mux.HandleFunc("GET /login", loginView)              			// Returns full page, login view (initial page)    			
	mux.HandleFunc("GET /logout", logoutView)       				// Returns HTML fragment (HTMX), logout view
	mux.HandleFunc("POST /login", loginPostView)       			// Returns HTML fragment (HTMX), login view

	//student role
	mux.Handle("GET /studentDashboard", requireRole("student", http.HandlerFunc(studentView)))			// Returns full page, student dashboard view
	mux.HandleFunc("GET /shop", shopView)        		 			// Returns full page, shop view
	mux.HandleFunc("GET /avatar", avatarView)        	 			// Returns full page, avatar view
	
	//teacher role
	mux.Handle("GET /teacherDashboard", requireRole("teacher", http.HandlerFunc(teacherView)))			// Returns full page, teacher dashboard view
	mux.HandleFunc("GET /teacherDashboard/edit", teacherEditView) 	// Returns HTML fragment (HTMX), edit view
	//admin role
	mux.Handle("GET /adminDashboard", requireRole("admin", http.HandlerFunc(adminView)))				// Returns full page, admin dashboard view
	mux.HandleFunc("GET /adminDashboard/edit", adminEditView)   	// Returns HTML fragment (HTMX), edit view

	log.Print("starting server on http://localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
