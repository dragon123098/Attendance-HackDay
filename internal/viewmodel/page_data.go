package viewmodel
 
import "github.com/dragon123098/Attendance-HackDay.git/internal/domain"
 
type TemplateData struct {
	User           *domain.User
	ErrorMessage   string
	SuccessMessage string
	Schedule       *domain.Schedule
	Schedules      map[string]*domain.Schedule
}
 
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
 
type PageData struct {
	Title              string
	Username           string
	AvatarImage        string
	Error              string
	Coins              int
	HeaderTitle        string
	HeaderSubtitle     string
	HeaderBadge        string
	AttendanceStatus   string
	AttendanceMessage  string
	CanMarkAttendance  bool
	WeeklySchedule     []ScheduleItemView
	UpcomingDoubleDays []DoubleDayView
	ActiveNav          string
	UseStudentCSS      bool
	ShopItems          []ShopItemView
	OwnedShopItems     []ShopItemView
	ShopMessage        string
}