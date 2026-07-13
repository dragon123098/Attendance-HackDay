package web

import (
	"context"
	"time"

	"github.com/dragon123098/Attendance-HackDay.git/internal/domain"
)

type StudentStore interface {
	LoadStudentState(context.Context, string) (domain.StudentState, error)
	MarkAttendanceAndAwardCoins(context.Context, string, string, string, int, time.Time) error
	PurchaseShopItem(context.Context, string, string, time.Time) error
	SaveAvatarConfig(context.Context, string, domain.AvatarConfig) error
}

var studentStore StudentStore
