package domain

import (
	"new_todo_project/pb"
)

type Todo struct {
	Id          int    `sql-col:"id"`
	Title       string `sql-col:"title"`
	Description string `sq--col:"description"`
	Done        bool   `sql-col:"done"`
}

func NewTodo(title, description string, done bool) *Todo {
	return &Todo{
		Title:       title,
		Description: description,
		Done:        done,
	}
}

func (t *Todo) MarkDone() {
	t.Done = true
}

func ToPbUser(u *User) *pb.User {
	var todoPbs []*pb.Todo
	for _, todo := range u.Todos {
		todoPbs = append(todoPbs, ToPbTodo(todo))
	}
	return &pb.User{
		Name:     u.Name,
		Password: u.Password,
		Email:    u.Email,
		TodoList: todoPbs,
	}
}

func ToPbTodo(t *Todo) *pb.Todo {
	return &pb.Todo{
		Id:   int64(t.Id),
		Name: t.Title,
		Done: t.Done,
	}
}

func ToDomainTodo(t *pb.TodoRequest) *Todo {
	return &Todo{
		Title:       t.Todo.Name,
		Description: t.Todo.Description,
		Done:        t.Todo.Done,
	}
}
