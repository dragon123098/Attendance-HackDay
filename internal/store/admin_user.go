package store

import (
	"context"

	"github.com/dragon123098/Attendance-HackDay.git/internal/domain"
)

// ListUsers returns the users shown on the admin user settings page.
func (s *SQLStore) ListUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			UserID,
			Name,
			Role,
			Email,
			COALESCE(ClassroomID, '')
		FROM Users
		ORDER BY Name, UserID;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Role, &user.Email, &user.ClassroomID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUserRole changes one user's role from the admin user settings page.
func (s *SQLStore) UpdateUserRole(ctx context.Context, userID string, role string) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE Users
		SET Role = $2
		WHERE UserID = $1;
	`, userID, role)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
