package main

import (
	"fmt"
	"net/http"
	"time"
)

const startingStudentCoins = 10
const attendanceRewardCoins = 1

func studentView(w http.ResponseWriter, r *http.Request) {
	user, ok := currentSessionUser(w, r)
	if !ok {
		return
	}

	data := PageData{
		Title:       "Student Dashboard",
		Username:    user.Name,
		AvatarImage: "/static/images/geraldIcon3.png",
		Coins:       getCoinBalance(user.UserID),
	}

	renderStudent(w, "studentDash.html", data)
}

func shopView(w http.ResponseWriter, r *http.Request) {
	user, ok := currentSessionUser(w, r)
	if !ok {
		return
	}

	data := PageData{
		Title:       "Shop",
		Username:    user.Name,
		AvatarImage: "/static/images/geraldIcon3.png",
		Coins:       getCoinBalance(user.UserID),
	}

	renderStudent(w, "shopView.html", data)
}

func avatarView(w http.ResponseWriter, r *http.Request) {
	user, ok := currentSessionUser(w, r)
	if !ok {
		return
	}

	data := PageData{
		Title:       "Avatar",
		Username:    user.Name,
		AvatarImage: "/static/images/geraldIcon3.png",
		Coins:       getCoinBalance(user.UserID),
	}

	renderStudent(w, "avatarView.html", data)
}

func shopBuyView(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented yet", http.StatusNotImplemented)
}

func attendanceView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := currentSessionUser(w, r)
	if !ok {
		return
	}

	today := time.Now().Format("2006-01-02")
	reward := attendanceRewardCoins
	if isDoubleDay(user.ClassroomID, time.Now()) {
		reward *= 2
	}

	awarded, err := markAttendanceAndAwardCoins(user, today, reward)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if awarded {
		saveData()
	}

	http.Redirect(w, r, "/studentDashboard", http.StatusSeeOther)
}

func currentSessionUser(w http.ResponseWriter, r *http.Request) (*User, bool) {
	userID, err := getSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, false
	}

	user, ok := app.Users[userID]
	if !ok {
		clearSessionUser(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, false
	}

	return user, true
}

func getCoinBalance(userID string) int {
	total := startingStudentCoins
	for _, tx := range app.Transactions {
		if tx.UserID == userID {
			total += tx.Amount
		}
	}
	return total
}

func markAttendanceAndAwardCoins(user *User, date string, reward int) (bool, error) {
	for i := range app.Attendance {
		rec := &app.Attendance[i]
		if rec.UserID == user.UserID && rec.ClassroomID == user.ClassroomID {
			for _, presentDate := range rec.Present {
				if presentDate == date {
					return false, nil
				}
			}

			rec.Present = append(rec.Present, date)
			app.Transactions = append(app.Transactions, CoinTransaction{
				UserID:      user.UserID,
				Amount:      reward,
				Timestamp:   time.Now().Format(time.RFC3339),
				Description: fmt.Sprintf("Attendance reward for %s", date),
			})
			return true, nil
		}
	}

	app.Attendance = append(app.Attendance, AttendanceRecord{
		UserID:      user.UserID,
		ClassroomID: user.ClassroomID,
		Present:     []string{date},
		Absent:      []string{},
	})

	app.Transactions = append(app.Transactions, CoinTransaction{
		UserID:      user.UserID,
		Amount:      reward,
		Timestamp:   time.Now().Format(time.RFC3339),
		Description: fmt.Sprintf("Attendance reward for %s", date),
	})

	return true, nil
}

func isDoubleDay(classroomID string, now time.Time) bool {
	weekday := now.Weekday().String()

	for _, sched := range app.Schedule {
		if sched.ClassroomID == classroomID && sched.DoubleDay && sched.DayOfWeek == weekday {
			return true
		}
	}

	return false
}