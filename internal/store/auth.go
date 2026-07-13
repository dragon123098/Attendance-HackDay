package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dragon123098/Attendance-HackDay.git/internal/domain"
)

// ErrUserNotFound means no SQL user matched the requested credential lookup.
var ErrUserNotFound = errors.New("user not found")

// FindUserByEmail loads the credential and routing data needed by the login flow.
func (s *SQLStore) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.findUser(ctx, "Email", email)
}

// FindUserByID loads the current SQL user referenced by an authenticated session.
func (s *SQLStore) FindUserByID(ctx context.Context, userID string) (domain.User, error) {
	return s.findUser(ctx, "UserID", userID)
}

func (s *SQLStore) findUser(ctx context.Context, column, value string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT TOP (1)
			UserID,
			Name,
			Role,
			Email,
			PasswordHash,
			COALESCE(ClassroomID, N'')
		FROM dbo.Users
		WHERE ` + column + ` = @p1
		ORDER BY UserID;
	`
	err := s.db.QueryRowContext(ctx, query, value).Scan(
		&user.UserID,
		&user.Name,
		&user.Role,
		&user.Email,
		&user.PasswordHash,
		&user.ClassroomID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
