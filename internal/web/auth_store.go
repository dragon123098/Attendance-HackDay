package web

import (
	"context"

	"github.com/PeterGrunig/Attendance-HackDay/internal/domain"
)

type AuthStore interface {
	FindUserByEmail(context.Context, string) (domain.User, error)
	FindUserByID(context.Context, string) (domain.User, error)
}

var authStore AuthStore
