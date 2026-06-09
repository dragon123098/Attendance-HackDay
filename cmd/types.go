package main

// AppState holds the application state, including users and schedules.
type AppState struct {
	Users    map[string]*User `json:"users"`
}

type User struct {
	Name    string `json:"name"`
	Role    string `json:"role"` // "student", "teacher", "admin"
	Email   string `json:"email"`
	PasswordHash string `json:"password_hash"`
	ClassroomID string `json:"classroom_id"` // for students, which classroom they belong to
}

type Classroom struct {
	Name string `json:"name"`
	ID  string `json:"id"`

}

type Schedule struct {
	ClassroomID string `json:"classroom_id"`
	DayOfWeek  string `json:"day_of_week"` // "Monday", "Tuesday", etc.
	StartTime   string `json:"start_time"`  // "09:00"
	EndTime     string `json:"end_time"`    // "10:00"
}