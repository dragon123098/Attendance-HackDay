package store

import (
	"context"
	"log"
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


//This will store a new Classroom in the database. It will be called from the admin tools.
func (s *SQLStore) CreateClassroom(ctx context.Context, classroom domain.Classroom) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dbo.Classrooms (ID, Name)
		VALUES (@p1, @p2);
	`, classroom.ID, classroom.Name); err != nil {
		log.Println("Error inserting classroom:", err)
		return err
	}

	return tx.Commit()
}