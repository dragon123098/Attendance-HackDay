package store

import (
	"context"

	"github.com/PeterGrunig/Attendance-HackDay/internal/domain"
)



func (s *SQLStore) TeacherListClassrooms(ctx context.Context, teacherID string) ([]domain.Classroom, error) {
		rows, err := s.db.QueryContext(ctx, `
		SELECT c.ID, c.Name, COALESCE(c.TeacherID, ''), COALESCE(cs.StudentID, '')
		FROM Classrooms AS c
		LEFT JOIN ClassroomStudents AS cs
			ON cs.ClassroomID = c.ID
		WHERE c.TeacherID = $1
		ORDER BY c.ID, cs.StudentID;
	`, teacherID)
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
