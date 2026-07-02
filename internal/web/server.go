package web

import (
	"io/fs"
	"net/http"

	"github.com/dragon123098/Attendance-HackDay.git/internal/view"
)

func NewRouter(studentStore AdminStudentStore) http.Handler {
	adminStudentStore = studentStore
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(view.FS, "static")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// auth routes
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("POST /logout", logoutView)

	// student routes
	mux.Handle("GET /studentDashboard", RequireRole("student", http.HandlerFunc(studentView)))
	mux.Handle("GET /shop", RequireRole("student", http.HandlerFunc(shopView)))
	mux.Handle("POST /shop/buy", RequireRole("student", http.HandlerFunc(shopBuyView)))
	mux.Handle("GET /avatar", RequireRole("student", http.HandlerFunc(avatarView)))
	mux.Handle("POST /avatar/preview", RequireRole("student", http.HandlerFunc(avatarPreviewView)))
	mux.Handle("POST /avatar", RequireRole("student", http.HandlerFunc(avatarSaveView)))
	mux.Handle("POST /attendance", RequireRole("student", http.HandlerFunc(attendanceView)))

	// teacher routes
	mux.Handle("GET /teacherDashboard", RequireRole("teacher", http.HandlerFunc(teacherView)))
	mux.HandleFunc("GET /teacherDashboard/edit", teacherEditView)

	// admin routes
	mux.Handle("GET /adminDashboard", RequireRole("admin", http.HandlerFunc(adminView)))
	mux.HandleFunc("GET /adminDashboard/edit", adminEditView)
	mux.HandleFunc("POST /classrooms", createClassroomView)
	mux.HandleFunc("GET /classrooms", listClassroomsView)
	mux.HandleFunc("GET /classrooms/edit", editClassrooms)
	mux.HandleFunc("POST /classrooms/edit", saveClassrooms)
	mux.HandleFunc("GET /addTeacher", createTeacher)
	mux.HandleFunc("POST /addTeacher", teacherCreateSubmitView)
	mux.Handle("GET /addStudent", RequireRole("admin", http.HandlerFunc(createStudent)))
	mux.Handle("POST /addStudent", RequireRole("admin", http.HandlerFunc(studentCreateSubmitView)))
	mux.HandleFunc("GET /userSettings", userSettingsView)
	mux.HandleFunc("POST /userSettings/role", updateUserRoleView)

	return mux
}
