package web

import (
	"io/fs"
	"net/http"
	"path"

	"github.com/dragon123098/Attendance-HackDay.git/internal/view"
)

type AppStore interface {
	AdminStudentStore
	AdminTeacherStore
	AdminClassroomStore
	AdminUserStore
	AuthStore
	StudentStore
}

func NewRouter(appStore AppStore) http.Handler {
	adminStudentStore = appStore
	adminTeacherStore = appStore
	adminClassroomStore = appStore
	adminUserStore = appStore
	authStore = appStore
	studentStore = appStore
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(view.FS, "static")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(staticFS))
	staticHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if extension := path.Ext(r.URL.Path); extension == ".css" || extension == ".js" {
			w.Header().Set("Cache-Control", "no-cache")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		fileServer.ServeHTTP(w, r)
	})
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// auth routes
	mux.HandleFunc("/", loginHandler)
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
	mux.Handle("POST /teacherDashboard/edit", RequireRole("teacher", http.HandlerFunc(teacherEditView)))

	// admin routes
	mux.Handle("GET /adminDashboard", RequireRole("admin", http.HandlerFunc(adminView)))
	mux.Handle("POST /adminDashboard/edit", RequireRole("admin", http.HandlerFunc(adminEditView)))
	mux.Handle("POST /classrooms", RequireRole("admin", http.HandlerFunc(createClassroomView)))
	mux.Handle("GET /classrooms", RequireRole("admin", http.HandlerFunc(listClassroomsView)))
	mux.Handle("GET /classrooms/edit", RequireRole("admin", http.HandlerFunc(editClassrooms)))
	mux.Handle("POST /classrooms/edit", RequireRole("admin", http.HandlerFunc(saveClassrooms)))
	mux.Handle("GET /addTeacher", RequireRole("admin", http.HandlerFunc(createTeacher)))
	mux.Handle("POST /addTeacher", RequireRole("admin", http.HandlerFunc(teacherCreateSubmitView)))
	mux.Handle("GET /addStudent", RequireRole("admin", http.HandlerFunc(createStudent)))
	mux.Handle("POST /addStudent", RequireRole("admin", http.HandlerFunc(studentCreateSubmitView)))
	mux.Handle("GET /userSettings", RequireRole("admin", http.HandlerFunc(userSettingsView)))
	mux.Handle("POST /userSettings/role", RequireRole("admin", http.HandlerFunc(updateUserRoleView)))

	return mux
}
