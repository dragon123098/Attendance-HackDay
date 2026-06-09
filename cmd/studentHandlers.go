package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const startingStudentCoins = 10
const attendanceRewardCoins = 1
const defaultAvatarImage = "/static/images/geraldIcon3.png"

type ScheduleItemView struct {
	DayOfWeek string
	StartTime string
	EndTime   string
	DoubleDay bool
	IsToday   bool
}

type DoubleDayView struct {
	DayOfWeek string
	StartTime string
	EndTime   string
}

type ShopItemView struct {
	ID          string
	Name        string
	Description string
	Price       int
	Owned       bool
	Image       string `json:"image"`
}

func studentView(w http.ResponseWriter, r *http.Request) {
	user, ok := currentSessionUser(w, r)
	if !ok {
		return
	}

	attendanceStatus, attendanceMessage, canMark := getTodayAttendanceState(user)
	weeklySchedule := getWeeklySchedule(user)
	upcomingDoubleDays := getUpcomingDoubleDays(user)

	data := PageData{
		Title:              "Student Dashboard",
		Username:           user.Name,
		AvatarImage:        getAvatarImage(user),
		Coins:              getCoinBalance(user.UserID),
		AttendanceStatus:   attendanceStatus,
		AttendanceMessage:  attendanceMessage,
		CanMarkAttendance:  canMark,
		WeeklySchedule:     weeklySchedule,
		UpcomingDoubleDays: upcomingDoubleDays,
		ActiveNav:          "home",
		UseStudentCSS: 		true,
	}

	renderStudent(w, "studentDash.html", data)
}

func avatarView(w http.ResponseWriter, r *http.Request) {
	user, ok := currentSessionUser(w, r)
	if !ok {
		return
	}

	data := PageData{
		Title:       "Avatar",
		Username:    user.Name,
		AvatarImage: getAvatarImage(user),
		Coins:       getCoinBalance(user.UserID),
		ActiveNav:   "avatar",
		UseStudentCSS: true,
	}

	renderStudent(w, "avatarView.html", data)
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

	now := time.Now()
	today := now.Format("2006-01-02")

	reward := attendanceRewardCoins
	if isDoubleDay(user.ClassroomID, now) {
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

func getAvatarImage(user *User) string {
	return defaultAvatarImage
}

func getTodayAttendanceState(user *User) (status string, message string, canMark bool) {
	today := time.Now().Format("2006-01-02")

	for _, rec := range app.Attendance {
		if rec.UserID == user.UserID && rec.ClassroomID == user.ClassroomID {
			for _, presentDate := range rec.Present {
				if presentDate == today {
					return "Present today", "Attendance already marked for today.", false
				}
			}
			return "Not marked yet", "Tap Mark Attendance to earn coins.", true
		}
	}

	return "Not marked yet", "Tap Mark Attendance to earn coins.", true
}

func getWeeklySchedule(user *User) []ScheduleItemView {
	items := make([]ScheduleItemView, 0)

	today := time.Now().Weekday().String()
	for _, sched := range app.Schedule {
		if sched.ClassroomID != user.ClassroomID {
			continue
		}
		items = append(items, ScheduleItemView{
			DayOfWeek: sched.DayOfWeek,
			StartTime: sched.StartTime,
			EndTime:   sched.EndTime,
			DoubleDay: sched.DoubleDay,
			IsToday:   sched.DayOfWeek == today,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return weekdayIndex(items[i].DayOfWeek) < weekdayIndex(items[j].DayOfWeek)
	})

	return items
}

func getUpcomingDoubleDays(user *User) []DoubleDayView {
	items := make([]DoubleDayView, 0)

	for _, sched := range app.Schedule {
		if sched.ClassroomID != user.ClassroomID || !sched.DoubleDay {
			continue
		}
		items = append(items, DoubleDayView{
			DayOfWeek: sched.DayOfWeek,
			StartTime: sched.StartTime,
			EndTime:   sched.EndTime,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return weekdayIndex(items[i].DayOfWeek) < weekdayIndex(items[j].DayOfWeek)
	})

	return items
}

func weekdayIndex(day string) int {
	switch day {
	case "Sunday":
		return 0
	case "Monday":
		return 1
	case "Tuesday":
		return 2
	case "Wednesday":
		return 3
	case "Thursday":
		return 4
	case "Friday":
		return 5
	case "Saturday":
		return 6
	default:
		return 7
	}
}

func shopView(w http.ResponseWriter, r *http.Request) {
	user, ok := currentSessionUser(w, r)
	if !ok {
		return
	}

	ensureShopState()
	seedShopItems()

	allItems, ownedItems := getShopItemViews(user.UserID)

	data := PageData{
		Title:          "Shop",
		Username:       user.Name,
		AvatarImage:    getAvatarImage(user),
		Coins:          getCoinBalance(user.UserID),
		ShopItems:      allItems,
		OwnedShopItems: ownedItems,
		ShopMessage:    r.URL.Query().Get("msg"),
		ActiveNav: "shop",
		UseStudentCSS: true,
	}

	renderStudent(w, "shopView.html", data)
}

func shopBuyView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := currentSessionUser(w, r)
	if !ok {
		return
	}

	ensureShopState()
	seedShopItems()

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/shop?msg="+url.QueryEscape("Invalid form submission."), http.StatusSeeOther)
		return
	}

	itemID := strings.TrimSpace(r.FormValue("item_id"))
	item, ok := app.ShopItems[itemID]
	if !ok {
		http.Redirect(w, r, "/shop?msg="+url.QueryEscape("That item does not exist."), http.StatusSeeOther)
		return
	}

	if userOwnsShopItem(user.UserID, itemID) {
		http.Redirect(w, r, "/shop?msg="+url.QueryEscape("You already own that item."), http.StatusSeeOther)
		return
	}

	balance := getCoinBalance(user.UserID)
	if balance < item.Price {
		http.Redirect(w, r, "/shop?msg="+url.QueryEscape("You do not have enough coins."), http.StatusSeeOther)
		return
	}

	app.Transactions = append(app.Transactions, CoinTransaction{
		UserID:      user.UserID,
		Amount:      -item.Price,
		Timestamp:   time.Now().Format(time.RFC3339),
		Description: fmt.Sprintf("Purchased %s", item.Name),
	})

	app.OwnedShopItems[user.UserID] = appendUniqueString(app.OwnedShopItems[user.UserID], itemID)
	saveData()

	http.Redirect(w, r, "/shop?msg="+url.QueryEscape("Purchase complete."), http.StatusSeeOther)
}

func ensureShopState() {
	if app.ShopItems == nil {
		app.ShopItems = map[string]*ShopItem{}
	}
	if app.OwnedShopItems == nil {
		app.OwnedShopItems = map[string][]string{}
	}
}

func seedShopItems() {
	if len(app.ShopItems) > 0 {
		return
	}

	app.ShopItems["hat_star"] = &ShopItem{
		ID:          "hat_star",
		Name:        "Star Hat",
		Price:       5,
		Description: "A bright hat for a standout student.",
	}
	app.ShopItems["trail_rainbow"] = &ShopItem{
		ID:          "trail_rainbow",
		Name:        "Rainbow Trail",
		Price:       8,
		Description: "A colorful trail effect for your avatar.",
	}
	app.ShopItems["cape_gold"] = &ShopItem{
		ID:          "cape_gold",
		Name:        "Golden Cape",
		Price:       12,
		Description: "A shiny cape for extra style.",
	}
	app.ShopItems["glasses_rocket"] = &ShopItem{
		ID:          "glasses_rocket",
		Name:        "Rocket Glasses",
		Price:       10,
		Description: "A bold accessory for your avatar.",
	}

	saveData()
}

func getShopItemViews(userID string) ([]ShopItemView, []ShopItemView) {
	items := make([]ShopItemView, 0, len(app.ShopItems))
	owned := make([]ShopItemView, 0)

	ids := make([]string, 0, len(app.ShopItems))
	for id := range app.ShopItems {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		item := app.ShopItems[id]
		view := ShopItemView{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Owned:       userOwnsShopItem(userID, item.ID),
		}
		items = append(items, view)
		if view.Owned {
			owned = append(owned, view)
		}
	}

	return items, owned
}

func userOwnsShopItem(userID, itemID string) bool {
	for _, ownedID := range app.OwnedShopItems[userID] {
		if ownedID == itemID {
			return true
		}
	}
	return false
}

func appendUniqueString(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}