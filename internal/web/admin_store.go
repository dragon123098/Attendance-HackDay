package web

import "context"

type AdminStudentStore interface {
	ListClassrooms(context.Context) ([]Classroom, error)
	CreateStudent(context.Context, User) error
}

type AdminTeacherStore interface {
	CreateTeacher(context.Context, User) error
}
