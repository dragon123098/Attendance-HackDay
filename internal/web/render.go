package web

import (
	"html/template"
	"net/http"
	"path/filepath"

)

//These two functions load pages for unauthenticated users.
func loadUnAuthTemplates(page string) (*template.Template, error) {
	return template.ParseFiles(
		filepath.Join("templates", "UnAuthBase.html"),
		filepath.Join("templates", page),
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


//Load Student Templates

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
	return template.ParseFiles(
		filepath.Join("templates", "Studentbase.html"),
		filepath.Join("templates", "partials", "topbar.html"),
		filepath.Join("templates", "partials", "StudentNavbar.html"),
		filepath.Join("templates", "partials", "footer.html"),
		filepath.Join("templates", page),
	)
}

func loadTeacherTemplates(page string) (*template.Template, error) {
	return template.ParseFiles(
		filepath.Join("templates", "adminBase.html"),
		filepath.Join("templates", "partials", "teacherNavBar.html"),
		filepath.Join("templates", "partials", "teacherHeader.html"),
		
		filepath.Join("templates", page),
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

func loadAdminTemplates (page string) (*template.Template, error) {
	return template.ParseFiles(
		filepath.Join("templates", "adminBase.html"),
		filepath.Join("templates", "partials", "adminHeader.html"),
		filepath.Join("templates", "partials", "adminNavBar.html"),
		
		filepath.Join("templates", page),
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