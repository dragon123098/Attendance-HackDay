package web

import "context"

type AdminStudentStore interface {
	ListClassrooms(context.Context) ([]Classroom, error)
	CreateStudent(context.Context, User) error
}

type AdminTeacherStore interface {
	CreateTeacher(context.Context, User) error
}

type AdminClassroomStore interface {
	ListClassrooms(context.Context) ([]Classroom, error)
	ListClassroomUsers(context.Context) (map[string]User, error)
	CreateClassroom(context.Context, Classroom) error
	UpdateClassroom(context.Context, string, Classroom) error
}
