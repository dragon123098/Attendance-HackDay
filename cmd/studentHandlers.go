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
const defaultAvatarBase = "gerald_icon"
const defaultAvatarImage = "/static/images/geraldIcon3.png"

type avatarBaseOption struct {
	ID    string
	Name  string
	Image string
}

type avatarCosmeticOption struct {
	ID        string
	SlotLabel string
	FormName  string
}

var avatarBaseOptions = []avatarBaseOption{
	{ID: "gerald_icon", Name: "Gerald Badge", Image: "/static/images/geraldIcon3.png"},
	{ID: "gerald_classic", Name: "Gerald Classic", Image: "/static/images/gerald.png"},
	{ID: "gerald_focus", Name: "Gerald Focus", Image: "/static/images/geraldIcon2.png"},
	{ID: "gerald_hero", Name: "Gerald Hero", Image: "/static/images/geraldIcon.png"},
}

var avatarCosmeticOptions = []avatarCosmeticOption{
	{ID: "hat_star", SlotLabel: "Headwear", FormName: "hair_style"},
	{ID: "cape_gold", SlotLabel: "Clothing", FormName: "clothing"},
	{ID: "glasses_rocket", SlotLabel: "Accessory", FormName: "accessory"},
	{ID: "trail_rainbow", SlotLabel: "Effect", FormName: "effect"},
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
		AvatarBadges:       getAvatarBadges(user.UserID),
		ActiveNav:          "home",
		UseStudentCSS:      true,
	}

	renderStudent(w, "studentDash.html", data)
}

func avatarView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := currentSessionUser(w, r)
	if !ok {
		return
	}

	ensureShopState()
	seedShopItems()
	ensureAvatarState()

	if r.Method == http.MethodPost {
		saveAvatarView(w, r, user)
		return
	}

	renderStudent(w, "avatarView.html", avatarPageData(user, r.URL.Query().Get("msg")))
}

func saveAvatarView(w http.ResponseWriter, r *http.Request, user *User) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/avatar?msg="+url.QueryEscape("Invalid form submission."), http.StatusSeeOther)
		return
	}

	config := getAvatarConfig(user.UserID)
	base := strings.TrimSpace(r.FormValue("base"))
	if base == "" {
		base = defaultAvatarBase
	}
	if !avatarBaseExists(base) {
		http.Redirect(w, r, "/avatar?msg="+url.QueryEscape("That avatar base does not exist."), http.StatusSeeOther)
		return
	}

	config.Base = base

	handledSlots := map[string]bool{}
	for _, option := range avatarCosmeticOptions {
		if handledSlots[option.FormName] {
			continue
		}
		handledSlots[option.FormName] = true

		selectedID := strings.TrimSpace(r.FormValue(option.FormName))
		if err := validateAvatarCosmeticSelection(user.UserID, option.FormName, selectedID); err != nil {
			http.Redirect(w, r, "/avatar?msg="+url.QueryEscape(err.Error()), http.StatusSeeOther)
			return
		}
		setAvatarSlot(&config, option.FormName, selectedID)
	}

	app.AvatarConfigs[user.UserID] = &config
	saveData()

	http.Redirect(w, r, "/avatar?msg="+url.QueryEscape("Avatar saved."), http.StatusSeeOther)
}

func avatarPageData(user *User, message string) PageData {
	data := PageData{
		Title:           "Avatar",
		Username:        user.Name,
		AvatarImage:     getAvatarImage(user),
		Coins:           getCoinBalance(user.UserID),
		AvatarBases:     getAvatarBaseViews(user.UserID),
		AvatarCosmetics: getAvatarCosmeticGroups(user.UserID),
		AvatarBadges:    getAvatarBadges(user.UserID),
		AvatarMessage:   message,
		ActiveNav:       "avatar",
		UseStudentCSS:   true,
	}

	return data
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
	config := getAvatarConfig(user.UserID)
	return avatarImageForBase(config.Base)
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
		AvatarBadges:   getAvatarBadges(user.UserID),
		ShopMessage:    r.URL.Query().Get("msg"),
		ActiveNav:      "shop",
		UseStudentCSS:  true,
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

func ensureAvatarState() {
	if app.AvatarConfigs == nil {
		app.AvatarConfigs = map[string]*AvatarConfig{}
	}
}

func seedShopItems() {
	changed := false
	for _, item := range defaultShopItems() {
		if _, exists := app.ShopItems[item.ID]; exists {
			continue
		}

		itemCopy := item
		app.ShopItems[item.ID] = &itemCopy
		changed = true
	}

	if changed {
		saveData()
	}
}

func defaultShopItems() []ShopItem {
	return []ShopItem{
		{
			ID:          "hat_star",
			Name:        "Star Hat",
			Price:       5,
			Description: "A bright hat for a standout student.",
		},
		{
			ID:          "trail_rainbow",
			Name:        "Rainbow Trail",
			Price:       8,
			Description: "A colorful trail effect for your avatar.",
		},
		{
			ID:          "cape_gold",
			Name:        "Golden Cape",
			Price:       12,
			Description: "A shiny cape for extra style.",
		},
		{
			ID:          "glasses_rocket",
			Name:        "Rocket Glasses",
			Price:       10,
			Description: "A bold accessory for your avatar.",
		},
	}
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

func getAvatarConfig(userID string) AvatarConfig {
	ensureAvatarState()

	config := AvatarConfig{Base: defaultAvatarBase}
	if saved, ok := app.AvatarConfigs[userID]; ok && saved != nil {
		config = *saved
	}
	if config.Base == "" || !avatarBaseExists(config.Base) {
		config.Base = defaultAvatarBase
	}

	return config
}

func avatarBaseExists(baseID string) bool {
	for _, option := range avatarBaseOptions {
		if option.ID == baseID {
			return true
		}
	}
	return false
}

func avatarImageForBase(baseID string) string {
	for _, option := range avatarBaseOptions {
		if option.ID == baseID {
			return option.Image
		}
	}
	return defaultAvatarImage
}

func getAvatarBaseViews(userID string) []AvatarBaseView {
	config := getAvatarConfig(userID)
	views := make([]AvatarBaseView, 0, len(avatarBaseOptions))

	for _, option := range avatarBaseOptions {
		views = append(views, AvatarBaseView{
			ID:       option.ID,
			Name:     option.Name,
			Image:    option.Image,
			Selected: option.ID == config.Base,
		})
	}

	return views
}

func getAvatarCosmeticGroups(userID string) []AvatarCosmeticGroupView {
	config := getAvatarConfig(userID)
	groupIndexByForm := map[string]int{}
	groups := make([]AvatarCosmeticGroupView, 0)

	for _, option := range avatarCosmeticOptions {
		groupIndex, exists := groupIndexByForm[option.FormName]
		if !exists {
			groups = append(groups, AvatarCosmeticGroupView{
				SlotLabel:  option.SlotLabel,
				FormName:   option.FormName,
				SelectedID: selectedAvatarSlot(config, option.FormName),
			})
			groupIndex = len(groups) - 1
			groupIndexByForm[option.FormName] = groupIndex
		}

		item := shopItemForID(option.ID)
		groups[groupIndex].Options = append(groups[groupIndex].Options, AvatarCosmeticView{
			ID:          option.ID,
			Name:        item.Name,
			Description: item.Description,
			SlotLabel:   option.SlotLabel,
			FormName:    option.FormName,
			Owned:       userOwnsShopItem(userID, option.ID),
			Selected:    selectedAvatarSlot(config, option.FormName) == option.ID,
		})
	}

	return groups
}

func validateAvatarCosmeticSelection(userID, formName, selectedID string) error {
	if selectedID == "" {
		return nil
	}

	for _, option := range avatarCosmeticOptions {
		if option.FormName == formName && option.ID == selectedID {
			if !userOwnsShopItem(userID, selectedID) {
				return fmt.Errorf("Buy %s in the shop before equipping it.", shopItemForID(selectedID).Name)
			}
			return nil
		}
	}

	return fmt.Errorf("That avatar cosmetic does not exist.")
}

func selectedAvatarSlot(config AvatarConfig, formName string) string {
	switch formName {
	case "hair_style":
		return config.HairStyle
	case "clothing":
		return config.Clothing
	case "accessory":
		return config.Accessory
	case "effect":
		return config.Effect
	default:
		return ""
	}
}

func setAvatarSlot(config *AvatarConfig, formName, selectedID string) {
	switch formName {
	case "hair_style":
		config.HairStyle = selectedID
	case "clothing":
		config.Clothing = selectedID
	case "accessory":
		config.Accessory = selectedID
	case "effect":
		config.Effect = selectedID
	}
}

func shopItemForID(itemID string) ShopItem {
	if item, ok := app.ShopItems[itemID]; ok && item != nil {
		return *item
	}

	for _, item := range defaultShopItems() {
		if item.ID == itemID {
			return item
		}
	}

	return ShopItem{
		ID:          itemID,
		Name:        itemID,
		Description: "Avatar cosmetic.",
	}
}

func getAvatarBadges(userID string) []AvatarBadgeView {
	config := getAvatarConfig(userID)
	badges := make([]AvatarBadgeView, 0)

	for _, option := range avatarCosmeticOptions {
		if selectedAvatarSlot(config, option.FormName) != option.ID {
			continue
		}

		item := shopItemForID(option.ID)
		badges = append(badges, AvatarBadgeView{
			Name:      item.Name,
			SlotLabel: option.SlotLabel,
		})
	}

	return badges
}
