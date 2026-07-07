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

// ListClassrooms returns classroom details and roster IDs used by admin forms.
func (s *SQLStore) ListClassrooms(ctx context.Context) ([]domain.Classroom, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.ID, c.Name, COALESCE(c.TeacherID, N''), COALESCE(cs.StudentID, N'')
		FROM dbo.Classrooms AS c
		LEFT JOIN dbo.ClassroomStudents AS cs
			ON cs.ClassroomID = c.ID
		ORDER BY c.ID, cs.StudentID;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	classrooms := []domain.Classroom{}
	classroomIndexes := map[string]int{}
	for rows.Next() {
		var (
			classroomID string
			name        string
			teacherID   string
			studentID   string
		)
		if err := rows.Scan(&classroomID, &name, &teacherID, &studentID); err != nil {
			return nil, err
		}

		index, ok := classroomIndexes[classroomID]
		if !ok {
			classrooms = append(classrooms, domain.Classroom{
				ID:        classroomID,
				Name:      name,
				TeacherID: teacherID,
			})
			index = len(classrooms) - 1
			classroomIndexes[classroomID] = index
		}

		if studentID != "" {
			classrooms[index].StudentIDs = append(classrooms[index].StudentIDs, studentID)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return classrooms, nil
}

// ListClassroomUsers returns users that can be displayed in classroom roster views.
func (s *SQLStore) ListClassroomUsers(ctx context.Context) (map[string]domain.User, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			UserID,
			Name,
			Role,
			Email,
			COALESCE(ClassroomID, N'')
		FROM dbo.Users
		ORDER BY UserID;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := map[string]domain.User{}
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Role, &user.Email, &user.ClassroomID); err != nil {
			return nil, err
		}
		users[user.UserID] = user
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
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
