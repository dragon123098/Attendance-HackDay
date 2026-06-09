package main

import (
	"fmt"
	"net/http"
	"time"
)

const attendanceRewardCoins = 10

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

	if user.Role != "student" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if user.ClassroomID == "" {
		http.Error(w, "student is not assigned to a classroom", http.StatusBadRequest)
		return
	}

	today := time.Now().Format("2006-01-02")

	recorded, err := markAttendanceAndAwardCoins(user, today)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if recorded {
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
	total := 0
	for _, tx := range app.Transactions {
		if tx.UserID == userID {
			total += tx.Amount
		}
	}
	return total
}

func markAttendanceAndAwardCoins(user *User, date string) (bool, error) {
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
				Amount:      attendanceRewardCoins,
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
		Amount:      attendanceRewardCoins,
		Timestamp:   time.Now().Format(time.RFC3339),
		Description: fmt.Sprintf("Attendance reward for %s", date),
	})

	return true, nil
}