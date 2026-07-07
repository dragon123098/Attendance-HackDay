package web

import (
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strings"

	datastore "github.com/dragon123098/Attendance-HackDay.git/internal/store"
	"golang.org/x/crypto/bcrypt"
)

var adminStudentStore AdminStudentStore
var adminTeacherStore AdminTeacherStore
var adminClassroomStore AdminClassroomStore

type ClassroomPageData struct {
	Title          string
	HeaderTitle    string
	HeaderSubtitle string
	HeaderBadge    string
	Classrooms     []*Classroom
}

type AdminDashboardPageData struct {
	Title          string
	Username       string
	HeaderTitle    string
	HeaderSubtitle string
	HeaderBadge    string
	Classrooms     []AdminClassroomView
}

type AdminClassroomView struct {
	Name     string
	ID       string
	Teacher  AdminClassroomPerson
	Students []AdminClassroomPerson
}

type AdminClassroomPerson struct {
	Name   string
	UserID string
}

type StudentCreatePageData struct {
	Title          string
	HeaderTitle    string
	HeaderSubtitle string
	HeaderBadge    string
	Classrooms     []ClassroomOption
}

type ClassroomOption struct {
	ID   string
	Name string
}

type UserSettingsPageData struct {
	Title          string
	HeaderTitle    string
	HeaderSubtitle string
	HeaderBadge    string
	Query          string
	Users          []UserSettingsRow
	RoleOptions    []RoleOption
}

type UserSettingsRow struct {
	Name   string
	UserID string
	Email  string
	Role   string
}

type RoleOption struct {
	Value string
	Label string
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

	data := AdminDashboardPageData{
		Title:          "Admin Dashboard",
		Username:       user.Name,
		HeaderTitle:    "Admin Dashboard",
		HeaderSubtitle: "Review classroom assignments and roster details.",
		HeaderBadge:    "Admin View",
		Classrooms:     buildAdminClassroomViews(),
	}

	renderAdmin(w, "adminDash.html", data)
}

func buildAdminClassroomViews() []AdminClassroomView {
	classroomIDs := make([]string, 0, len(app.Classrooms))
	for id := range app.Classrooms {
		classroomIDs = append(classroomIDs, id)
	}
	sort.Strings(classroomIDs)

	classrooms := make([]AdminClassroomView, 0, len(classroomIDs))
	for _, id := range classroomIDs {
		classroom := app.Classrooms[id]
		studentIDs := assignedStudentIDs(classroom)
		students := make([]AdminClassroomPerson, 0, len(studentIDs))

		for _, studentID := range studentIDs {
			students = append(students, adminClassroomPerson(studentID))
		}

		classrooms = append(classrooms, AdminClassroomView{
			Name:     classroom.Name,
			ID:       classroom.ID,
			Teacher:  adminClassroomPerson(classroom.TeacherID),
			Students: students,
		})
	}

	return classrooms
}

func assignedStudentIDs(classroom *Classroom) []string {
	seen := make(map[string]bool)
	for _, studentID := range classroom.StudentIDs {
		studentID = strings.TrimSpace(studentID)
		if studentID != "" {
			seen[studentID] = true
		}
	}

	for userID, user := range app.Users {
		if user.Role == "student" && user.ClassroomID == classroom.ID {
			seen[userID] = true
		}
	}

	studentIDs := make([]string, 0, len(seen))
	for studentID := range seen {
		studentIDs = append(studentIDs, studentID)
	}
	sort.Strings(studentIDs)

	return studentIDs
}

func adminClassroomPerson(userID string) AdminClassroomPerson {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return AdminClassroomPerson{Name: "Unassigned"}
	}

	user, ok := app.Users[userID]
	if !ok {
		return AdminClassroomPerson{
			Name:   "Unknown user",
			UserID: userID,
		}
	}

	return AdminClassroomPerson{
		Name:   user.Name,
		UserID: user.UserID,
	}
}

func classroomOptions() []ClassroomOption {
	classroomIDs := make([]string, 0, len(app.Classrooms))
	for id := range app.Classrooms {
		classroomIDs = append(classroomIDs, id)
	}
	sort.Strings(classroomIDs)

	options := make([]ClassroomOption, 0, len(classroomIDs))
	for _, id := range classroomIDs {
		classroom := app.Classrooms[id]
		options = append(options, ClassroomOption{
			ID:   classroom.ID,
			Name: classroom.Name,
		})
	}

	return options
}

func classroomOptionsFromStore(classrooms []Classroom) []ClassroomOption {
	sort.Slice(classrooms, func(i, j int) bool {
		return classrooms[i].ID < classrooms[j].ID
	})

	options := make([]ClassroomOption, 0, len(classrooms))
	for _, classroom := range classrooms {
		options = append(options, ClassroomOption{
			ID:   classroom.ID,
			Name: classroom.Name,
		})
	}

	return options
}

func userRoleOptions() []RoleOption {
	return []RoleOption{
		{Value: "student", Label: "Student"},
		{Value: "teacher", Label: "Teacher"},
		{Value: "admin", Label: "Admin"},
	}
}

func buildUserSettingsRows(query string) []UserSettingsRow {
	query = strings.ToLower(strings.TrimSpace(query))
	rows := make([]UserSettingsRow, 0, len(app.Users))

	for _, user := range app.Users {
		if query != "" && !userMatchesSearch(user, query) {
			continue
		}

		rows = append(rows, UserSettingsRow{
			Name:   user.Name,
			UserID: user.UserID,
			Email:  user.Email,
			Role:   user.Role,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		leftName := strings.ToLower(rows[i].Name)
		rightName := strings.ToLower(rows[j].Name)
		if leftName == rightName {
			return strings.ToLower(rows[i].UserID) < strings.ToLower(rows[j].UserID)
		}
		return leftName < rightName
	})

	return rows
}

func userMatchesSearch(user *User, query string) bool {
	values := []string{
		user.Name,
		user.Email,
		user.UserID,
		user.Role,
	}

	for _, value := range values {
		if strings.Contains(strings.ToLower(value), query) {
			return true
		}
	}

	return false
}

func isValidUserRole(role string) bool {
	for _, option := range userRoleOptions() {
		if role == option.Value {
			return true
		}
	}

	return false
}

func adminUserCount() int {
	count := 0
	for _, user := range app.Users {
		if user.Role == "admin" {
			count++
		}
	}
	return count
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

func userSettingsView(w http.ResponseWriter, r *http.Request) {
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

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	data := UserSettingsPageData{
		Title:          "User Settings",
		HeaderTitle:    "User Settings",
		HeaderSubtitle: "Search users and manage their roles.",
		HeaderBadge:    "Admin View",
		Query:          query,
		Users:          buildUserSettingsRows(query),
		RoleOptions:    userRoleOptions(),
	}

	renderAdmin(w, "userSettings.html", data)
}

func updateUserRoleView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUserID, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	currentUser, ok := app.Users[currentUserID]
	if !ok {
		clearSessionUser(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if currentUser.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}

	targetUserID := strings.TrimSpace(r.FormValue("user_id"))
	role := strings.TrimSpace(r.FormValue("role"))
	query := strings.TrimSpace(r.FormValue("q"))

	targetUser, ok := app.Users[targetUserID]
	if !ok {
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}

	if !isValidUserRole(role) {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	// Keep at least one reachable admin account after every role update.
	if targetUserID == currentUserID && role != "admin" {
		http.Error(w, "you cannot remove your own admin role", http.StatusBadRequest)
		return
	}

	if targetUser.Role == "admin" && role != "admin" && adminUserCount() <= 1 {
		http.Error(w, "cannot remove the last admin role", http.StatusBadRequest)
		return
	}

	targetUser.Role = role
	saveData()

	redirectTo := "/userSettings"
	if query != "" {
		redirectTo += "?q=" + url.QueryEscape(query)
	}
	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}

func listClassroomsView(w http.ResponseWriter, r *http.Request) {
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

func createClassroomView(w http.ResponseWriter, r *http.Request) {
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

	err := adminClassroomStore.CreateClassroom(r.Context(), *classroom)
	if err != nil {
		http.Error(w, "could not save classroom", http.StatusInternalServerError)
		return
	}

	
	http.Redirect(
		w,
		r,
		"/adminDashboard",
		http.StatusSeeOther,
	)
	
}

func editClassrooms(w http.ResponseWriter, r *http.Request) {
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
		Title:          "Classrooms",
		HeaderTitle:    "Admin Tools",
		HeaderSubtitle: "Manage classroom settings from here.",
		HeaderBadge:    "Admin View",
		Classrooms:     classrooms,
	}

	renderAdmin(w, "editClassrooms.html", data)
}

//This saves the edited classrooms. Right now it doesn't work how I want it to. 
func saveClassrooms(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}

	//originalID := r.FormValue("original_id")

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

	err := adminClassroomStore.CreateClassroom(r.Context(), *classroom)
	if err != nil {
		http.Error(w, "could not save classroom", http.StatusInternalServerError)
		return
	}
	

	http.Redirect(
		w,
		r,
		"/adminDashboard",
		http.StatusSeeOther,
	)
}

// MetaData for teacher information
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

	if adminTeacherStore == nil {
		http.Error(w, "teacher store is not configured", http.StatusInternalServerError)
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

	if err := adminTeacherStore.CreateTeacher(r.Context(), *teacher); err != nil {
		if errors.Is(err, datastore.ErrUserAlreadyExists) {
			http.Error(w, "teacher id already exists", http.StatusConflict)
			return
		}
		http.Error(w, "could not create teacher", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/adminDashboard", http.StatusSeeOther)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
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

	if adminStudentStore == nil {
		http.Error(w, "student store is not configured", http.StatusInternalServerError)
		return
	}

	classrooms, err := adminStudentStore.ListClassrooms(r.Context())
	if err != nil {
		http.Error(w, "could not load classrooms", http.StatusInternalServerError)
		return
	}

	data := StudentCreatePageData{
		Title:          "Add Student",
		HeaderTitle:    "Admin Tools",
		HeaderSubtitle: "Create a new student account.",
		HeaderBadge:    "Admin View",
		Classrooms:     classroomOptionsFromStore(classrooms),
	}

	renderAdmin(w, "createStudent.html", data)
}

func studentCreateSubmitView(w http.ResponseWriter, r *http.Request) {
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
	classroomID := strings.TrimSpace(r.FormValue("classroom_id"))
	password := r.FormValue("password")

	if userID == "" || name == "" || email == "" || password == "" || classroomID == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	if adminStudentStore == nil {
		http.Error(w, "student store is not configured", http.StatusInternalServerError)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "could not hash password", http.StatusInternalServerError)
		return
	}

	student := &User{
		Name:         name,
		Role:         "student",
		Email:        email,
		PasswordHash: string(hash),
		ClassroomID:  classroomID,
		UserID:       userID,
	}

	if err := adminStudentStore.CreateStudent(r.Context(), *student); err != nil {
		switch {
		case errors.Is(err, datastore.ErrClassroomNotFound):
			http.Error(w, "classroom does not exist", http.StatusBadRequest)
		case errors.Is(err, datastore.ErrUserAlreadyExists):
			http.Error(w, "student id already exists", http.StatusConflict)
		default:
			http.Error(w, "could not create student", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/adminDashboard", http.StatusSeeOther)
}
