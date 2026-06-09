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
	mux.Handle("GET /studentDashboard", RequireRole("student", http.HandlerFunc(studentView)))
	mux.Handle("GET /shop", RequireRole("student", http.HandlerFunc(shopView)))
	mux.Handle("POST /shop/buy", RequireRole("student", http.HandlerFunc(shopBuyView)))
	mux.Handle("GET /avatar", RequireRole("student", http.HandlerFunc(avatarView)))
	mux.Handle("POST /attendance", RequireRole("student", http.HandlerFunc(attendanceView)))

	// teacher routes
	mux.Handle("GET /teacherDashboard", RequireRole("teacher", http.HandlerFunc(teacherView)))
	mux.HandleFunc("GET /teacherDashboard/edit", teacherEditView)
	//mux.HandleFunc("GET /classroom", classroomView)
	//mux.HandleFunc("POST /classroom/:id/assign-student", teacherAssignStudentView)
	//mux.HandleFunc("POST /classroom/mark-attendance", markAttendanceView)

	// admin routes
	mux.Handle("GET /adminDashboard", RequireRole("admin", http.HandlerFunc(adminView)))
	mux.HandleFunc("GET /adminDashboard/edit", adminEditView)
	//mux.HandleFunc("POST /classrooms", createClassroomView)
	//mux.HandleFunc("GET /classrooms", listClassroomsView)
	//mux.HandleFunc("POST /classrooms/:id/assign-teacher", assignTeacherView)
	//mux.HandleFunc("POST /classrooms/:id/assign-student", assignStudentView)

	log.Print("starting server on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}