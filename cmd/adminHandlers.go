package main

import (
	"net/http"
	"strings"
)

func teacherView(w http.ResponseWriter, r *http.Request) {
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
		Title:    "Teacher Dashboard",
		Username: user.Name,
	}

	renderTeacher(w, "teacherDash.html", data)
}

func teacherEditView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	if user.Role != "teacher" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func adminView(w http.ResponseWriter, r *http.Request) {
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
		Title:    "Admin Dashboard",
		Username: user.Name,
	}

	renderAdmin(w, "adminDash.html", data)
}

func adminEditView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func listClassroomsView (w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	renderAdmin(w, "classrooms.html", nil)

}

func createClassroomView (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	id := r.FormValue("id")
	teacherID := r.FormValue("teacher_id")
	studentIDsRaw := r.FormValue("student_ids")

	var studentIDs []string
	for _, line := range strings.Split(studentIDsRaw, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			studentIDs = append(studentIDs, line)
		}
	}

	classroom := &Classroom{
		Name:       name,
		ID:         id,
		TeacherID:  teacherID,
		StudentIDs: studentIDs,
	}

	if app.Classrooms == nil {
		app.Classrooms = make(map[string]*Classroom)
	}

	app.Classrooms[id] = classroom

	saveData()
}