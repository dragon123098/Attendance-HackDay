package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dragon123098/Attendance-HackDay.git/internal/domain"
)

// ErrClassroomNotFound means a write referenced a classroom ID that is not in SQL Server.
var ErrClassroomNotFound = errors.New("classroom not found")

// ErrUserAlreadyExists means a requested user ID is already stored in SQL Server.
var ErrUserAlreadyExists = errors.New("user already exists")

// SQLStore owns SQL Server data access for flows that have moved off data.json.
type SQLStore struct {
	db *sql.DB
}

// NewSQLStore wraps an existing database handle without taking ownership of closing it.
func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

// ListClassrooms returns the classroom options used by admin forms.
func (s *SQLStore) ListClassrooms(ctx context.Context) ([]domain.Classroom, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT ID, Name, COALESCE(TeacherID, N'')
		FROM dbo.Classrooms
		ORDER BY ID;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	classrooms := []domain.Classroom{}
	for rows.Next() {
		var classroom domain.Classroom
		if err := rows.Scan(&classroom.ID, &classroom.Name, &classroom.TeacherID); err != nil {
			return nil, err
		}
		classrooms = append(classrooms, classroom)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return classrooms, nil
}

// CreateStudent inserts a student and classroom assignment in one SQL transaction.
func (s *SQLStore) CreateStudent(ctx context.Context, student domain.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := requireClassroom(ctx, tx, student.ClassroomID); err != nil {
		return err
	}
	if err := requireNewUser(ctx, tx, student.UserID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dbo.Users (UserID, Name, Role, Email, PasswordHash, ClassroomID)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6);
	`, student.UserID, student.Name, student.Role, student.Email, student.PasswordHash, student.ClassroomID); err != nil {
		return err
	}

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
	`, student.ClassroomID, student.UserID); err != nil {
		return err
	}

	return tx.Commit()
}

func requireClassroom(ctx context.Context, tx *sql.Tx, classroomID string) error {
	var exists int
	err := tx.QueryRowContext(ctx, `
		SELECT 1
		FROM dbo.Classrooms
		WHERE ID = @p1;
	`, classroomID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrClassroomNotFound
	}
	return err
}

func requireNewUser(ctx context.Context, tx *sql.Tx, userID string) error {
	var exists int
	err := tx.QueryRowContext(ctx, `
		SELECT 1
		FROM dbo.Users
		WHERE UserID = @p1;
	`, userID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	return ErrUserAlreadyExists
}
