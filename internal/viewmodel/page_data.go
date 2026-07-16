package viewmodel

import "github.com/dragon123098/Attendance-HackDay.git/internal/domain"

type TemplateData struct {
	User           *domain.User
	ErrorMessage   string
	SuccessMessage string
	Schedule       *domain.Schedule
	Schedules      map[string]*domain.Schedule
}

type WeeklyAssignmentView struct {
	Subject string
	Title   string
	DueTime string
}

type WeeklyScheduleDayView struct {
	DayName     string
	DateLabel   string
	DateISO     string
	IsToday     bool
	Assignments []WeeklyAssignmentView
}

type DoubleDayView struct {
	DayOfWeek string
	StartTime string
	EndTime   string
}

type ShopItemView struct {
	ID                string
	Name              string
	Description       string
	Price             int
	Owned             bool
	Slot              string
	Image             string `json:"image"`
	ThemeBackgroundID string
}

type ThemeBackgroundOptionView struct {
	ID    string
	Label string
}

type AvatarBaseOptionView struct {
	ID       string
	Label    string
	Image    string
	Selected bool
}

type AvatarCosmeticOptionView struct {
	ID       string
	Label    string
	Slot     string
	Image    string
	Owned    bool
	Selected bool
}

type AvatarLayerView struct {
	ID    string
	Label string
	Slot  string
	Image string
}

type AvatarPreviewView struct {
	BaseLabel      string
	BaseImage      string
	HairStyleLabel string
	ClothingLabel  string
	AccessoryLabel string
	EffectLabel    string
	HasCosmetics   bool
	Layers         []AvatarLayerView
}

type PageData struct {
	Title                  string
	Username               string
	AvatarImage            string
	AvatarSummary          []string
	Error                  string
	Coins                  int
	HeaderTitle            string
	HeaderSubtitle         string
	HeaderBadge            string
	AttendanceStatus       string
	AttendanceMessage      string
	CanMarkAttendance      bool
	CurrentWeekLabel       string
	WeeklyAssignmentDays   []WeeklyScheduleDayView
	UpcomingDoubleDays     []DoubleDayView
	ActiveNav              string
	UseStudentCSS          bool
	AvatarShopItems        []ShopItemView
	BackgroundShopItems    []ShopItemView
	OwnedShopItems         []ShopItemView
	ThemeBackgroundOptions []ThemeBackgroundOptionView
	ShopMessage            string
	AvatarBaseOptions      []AvatarBaseOptionView
	AvatarHairOptions      []AvatarCosmeticOptionView
	AvatarClothOptions     []AvatarCosmeticOptionView
	AvatarAccessOptions    []AvatarCosmeticOptionView
	AvatarEffectOptions    []AvatarCosmeticOptionView
	AvatarPreview          AvatarPreviewView
	AvatarMessage          string
	AvatarError            string
}
