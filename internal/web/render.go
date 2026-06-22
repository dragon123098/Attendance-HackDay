package web

import (
	"html/template"
	"io/fs"
	"net/http"
	"path"

	"github.com/dragon123098/Attendance-HackDay.git/internal/view"
)

func pageTemplate(page string) string {
	candidate := path.Join("pages", page)
	if _, err := fs.Stat(view.FS, candidate); err == nil {
		return candidate
	}

	return path.Join("pages", "popups", page)
}

// These two functions load pages for unauthenticated users.
func loadUnAuthTemplates(page string) (*template.Template, error) {
	return template.ParseFS(
		view.FS,
		"components/UnAuthBase.html",
		pageTemplate(page),
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

// Load Student Templates

func renderStudent(w http.ResponseWriter, page string, data any) {
	tmpl, err := loadStudentTemplates(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadStudentTemplates(page string) (*template.Template, error) {
	return template.ParseFS(
		view.FS,
		"components/Studentbase.html",
		"components/topbar.html",
		"components/StudentNavbar.html",
		"components/footer.html",
		pageTemplate(page),
	)
}

func loadTeacherTemplates(page string) (*template.Template, error) {
	return template.ParseFS(
		view.FS,
		"components/adminBase.html",
		"components/teacherNavBar.html",
		"components/teacherHeader.html",
		pageTemplate(page),
	)
}

func renderTeacher(w http.ResponseWriter, page string, data any) {
	tmpl, err := loadTeacherTemplates(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadAdminTemplates(page string) (*template.Template, error) {
	return template.ParseFS(
		view.FS,
		"components/adminBase.html",
		"components/adminHeader.html",
		"components/adminNavBar.html",
		pageTemplate(page),
	)
}

func renderAdmin(w http.ResponseWriter, page string, data any) {
	tmpl, err := loadAdminTemplates(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
