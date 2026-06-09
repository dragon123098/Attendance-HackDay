package main

type TemplateData struct {
	User           *User                // logged in user
	ErrorMessage   string               // error message to display
	SuccessMessage string               // success message to display
	Schedule       *Schedule            // current user's schedule
	Schedules      map[string]*Schedule // all schedules (for admin page)
}