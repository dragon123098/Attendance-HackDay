package store

import (
	"context"
	"database/sql"
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

// CreateClassroom stores a new classroom and its initial student assignments from the admin tools.
func (s *SQLStore) CreateClassroom(ctx context.Context, classroom domain.Classroom) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dbo.Classrooms (ID, Name, TeacherID)
		VALUES (@p1, @p2, @p3);
	`, classroom.ID, classroom.Name, classroom.TeacherID); err != nil {
		log.Println("Error inserting classroom:", err)
		return err
	}

	if err := insertClassroomStudents(ctx, tx, classroom.ID, classroom.StudentIDs); err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateClassroom saves edits to an existing classroom and replaces its student assignments.
func (s *SQLStore) UpdateClassroom(ctx context.Context, originalID string, classroom domain.Classroom) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE dbo.Classrooms
		SET ID = @p2, Name = @p3, TeacherID = @p4
		WHERE ID = @p1;
	`, originalID, classroom.ID, classroom.Name, classroom.TeacherID)
	if err != nil {
		log.Println("Error updating classroom:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrClassroomNotFound
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM dbo.ClassroomStudents
		WHERE ClassroomID = @p1 OR ClassroomID = @p2;
	`, originalID, classroom.ID); err != nil {
		log.Println("Error clearing classroom students:", err)
		return err
	}

	if err := insertClassroomStudents(ctx, tx, classroom.ID, classroom.StudentIDs); err != nil {
		return err
	}

	return tx.Commit()
}

func insertClassroomStudents(ctx context.Context, tx *sql.Tx, classroomID string, studentIDs []string) error {
	for _, studentID := range studentIDs {
		if _, err := tx.ExecContext(ctx, `
			IF NOT EXISTS (
				SELECT 1
				FROM dbo.ClassroomStudents
				WHERE ClassroomID = @p1 AND StudentID = @p2
			)
			BEGIN
				INSERT INTO dbo.ClassroomStudents (ClassroomID, StudentID)
				VALUES (@p1, @p2);
			END;
		`, classroomID, studentID); err != nil {
			log.Println("Error inserting classroom student:", err)
			return err
		}
	}

	return nil
}
