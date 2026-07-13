package web

import (
	"context"

	"github.com/dragon123098/Attendance-HackDay.git/internal/domain"
)

type AuthStore interface {
	FindUserByEmail(context.Context, string) (domain.User, error)
	FindUserByID(context.Context, string) (domain.User, error)
}

var authStore AuthStore
