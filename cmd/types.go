package main

// AppState holds the application state, including users and schedules.
type AppState struct {
    Users           map[string]*User           `json:"users"`
    Classrooms      map[string]*Classroom      `json:"classrooms"`
    ShopItems       map[string]*ShopItem       `json:"shop_items"`
    AvatarConfigs   map[string]*AvatarConfig   `json:"avatar_configs"`   // keyed by userID
    Transactions    []CoinTransaction          `json:"transactions"`
    Attendance      []AttendanceRecord         `json:"attendance"`
    Schedule        []Schedule                 `json:"schedule"`
}

type User struct {
	Name    string `json:"name"`
	Role    string `json:"role"` // "student", "teacher", "admin"
	Email   string `json:"email"`
	PasswordHash string `json:"password_hash"`
	ClassroomID string `json:"classroom_id"` // for students, which classroom they belong to
	UserID string `json:"user_id"`
}

type Classroom struct {
	Name string `json:"name"`
	ID  string `json:"id"`
	TeacherID string `json:"teacher_id"`
	StudentIDs []string `json:"student_ids"`
}

type Schedule struct {
	ClassroomID string `json:"classroom_id"`
	DayOfWeek  string `json:"day_of_week"` // "Monday", "Tuesday", etc.
	StartTime   string `json:"start_time"`  // "09:00"
	EndTime     string `json:"end_time"`    // "10:00"
	DoubleDay   bool   `json:"double_day"`  // true if this schedule is a double day (e.g. Monday and Wednesday)
}

type ShopItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"` // price in coins
	Description string `json:"description"`
}

type AvatarConfig struct {
	Base string `json:"base"`
	HairStyle string `json:"hair_style"`
	Clothing  string `json:"clothing"`
	Accessory string `json:"accessory"`
}

type CoinTransaction struct {
	UserID    string `json:"user_id"`
	Amount    int    `json:"amount"` // positive for earning coins, negative for spending coins
	Timestamp string `json:"timestamp"`
	Description string `json:"description"`
}

type AttendanceRecord struct {
	UserID    string `json:"user_id"`
	ClassroomID string `json:"classroom_id"`
	Present []string `json:"present"` // list of dates when the student was present
	Absent  []string `json:"absent"`  // list of dates when the student was absent	
}