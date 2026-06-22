package web

import (
	"net/http"
	"strings"
	"golang.org/x/crypto/bcrypt"
)

type ClassroomPageData struct {
	Title         string
	HeaderTitle    string
	HeaderSubtitle string
	HeaderBadge    string
	Classrooms     []*Classroom
}

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

func editClassrooms (w http.ResponseWriter, r *http.Request) {
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
	
	
	classrooms := make([]*Classroom, 0, len(app.Classrooms))

	for _, classroom := range app.Classrooms {
		classrooms = append(classrooms, classroom)
	}

	data := ClassroomPageData{
		Title:         "Classrooms",
		HeaderTitle:    "Admin Tools",
		HeaderSubtitle: "Manage classroom settings from here.",
		HeaderBadge:    "Admin View",
		Classrooms:     classrooms,
	}

	renderAdmin(w, "editClassrooms.html", data)
}

func saveClassrooms(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}

	originalID := r.FormValue("original_id")

	name := r.FormValue("name")
	id := r.FormValue("id")
	teacherID := r.FormValue("teacher_id")

	var studentIDs []string

	for _, line := range strings.Split(
		r.FormValue("student_ids"),
		"\n",
	) {
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

	// Handle ID changes
	if originalID != id {
		delete(app.Classrooms, originalID)
	}

	app.Classrooms[id] = classroom

	saveData()

	http.Redirect(
		w,
		r,
		"/adminDashboard",
		http.StatusSeeOther,
	)
}

//MetaData for teacher information
type TeacherCreatePageData struct {
	Title          string
	HeaderTitle    string
	HeaderSubtitle string
	HeaderBadge    string
}

func createTeacher(w http.ResponseWriter, r *http.Request) {
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
	data := TeacherCreatePageData{
		Title:          "Add Teacher",
		HeaderTitle:    "Admin Tools",
		HeaderSubtitle: "Create a new teacher account.",
		HeaderBadge:    "Admin View",
	}

	renderAdmin(w, "createTeacher.html", data)
}

func teacherCreateSubmitView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}

	userID := strings.TrimSpace(r.FormValue("user_id"))
	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if userID == "" || name == "" || email == "" || password == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	if app.Users == nil {
		app.Users = make(map[string]*User)
	}

	if _, exists := app.Users[userID]; exists {
		http.Error(w, "teacher id already exists", http.StatusConflict)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "could not hash password", http.StatusInternalServerError)
		return
	}

	teacher := &User{
		Name:         name,
		Role:         "teacher",
		Email:        email,
		PasswordHash: string(hash),
		ClassroomID:  "",
		UserID:       userID,
	}

	app.Users[userID] = teacher
	saveData()

	//http.Redirect(w, r, "/adminDashboard", http.StatusSeeOther)
}

