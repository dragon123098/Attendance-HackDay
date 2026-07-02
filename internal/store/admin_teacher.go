package store

import (
	"context"

	"github.com/dragon123098/Attendance-HackDay.git/internal/domain"
)

// CreateTeacher inserts a teacher account created from the admin tools.
func (s *SQLStore) CreateTeacher(ctx context.Context, teacher domain.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := requireNewUser(ctx, tx, teacher.UserID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dbo.Users (UserID, Name, Role, Email, PasswordHash, ClassroomID)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6);
	`, teacher.UserID, teacher.Name, teacher.Role, teacher.Email, teacher.PasswordHash, teacher.ClassroomID); err != nil {
		return err
	}

	return tx.Commit()
}
