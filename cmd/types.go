package main

// AppState holds the application state, including users and schedules.
type AppState struct {
	Users    map[string]*User
	Schedule map[string]*Schedule
}

type User struct {
}

type Classroom struct {
}

type Schedule struct {
}