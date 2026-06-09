package main

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

// PageData holds transient page-only values.
type PageData struct {
	Title             string
	Username          string
	AvatarImage       string
	Error             string
	Coins             int
	HeaderTitle       string
	HeaderSubtitle    string
	HeaderBadge       string
	AttendanceStatus  string
	AttendanceMessage string
	CanMarkAttendance bool
	WeeklySchedule    []ScheduleItemView
	UpcomingDoubleDays []DoubleDayView
	ActiveNav         string
	UseStudentCSS bool
	ShopItems       []ShopItemView
	OwnedShopItems  []ShopItemView
	ShopMessage     string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderUnAuth(w, "login.html", PageData{Title: "Login"})
		return

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form submission", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Email and password are required.",
			})
			return
		}

		var user *User
		for _, u := range app.Users {
			if u.Email == email {
				user = u
				break
			}
		}

		if user == nil {
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Invalid email or password.",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			renderUnAuth(w, "login.html", PageData{
				Title: "Login",
				Error: "Invalid email or password.",
			})
			return
		}

		if err := createSession(w, user.UserID); err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		switch user.Role {
		case "student":
			http.Redirect(w, r, "/studentDashboard", http.StatusSeeOther)
		case "teacher":
			http.Redirect(w, r, "/teacherDashboard", http.StatusSeeOther)
		case "admin":
			http.Redirect(w, r, "/adminDashboard", http.StatusSeeOther)
		default:
			http.Error(w, "invalid user role", http.StatusForbidden)
		}

		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func logoutView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clearSessionUser(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}