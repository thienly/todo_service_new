package domain

import (
	"errors"
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

type User struct {
	Id       int    `sql-col:"user_id"`
	Name     string `sql-col:"username"`
	Password string `sql-col:"password"`
	Email    string `sql-col:"email"`
	Todos    []*Todo
}

func NewUser(name, password, email string) *User {
	return &User{Name: name,
		Password: password,
		Email:    email}
}

func (u *User) ChangePassword(oldPassword, newPassword string) (bool, error) {
	if u.Password != oldPassword {
		return false, errors.New("Old password is not matched with new password")
	}
	u.Password = newPassword
	return true, nil
}

func (u *User) AddNewTodo(todo *Todo) error {
	if todo == nil {
		return errors.New("todo is not null")
	}
	u.Todos = append(u.Todos, todo)
	return nil
}

func (u *User) MarkDone(id int) error {
	found := false
	for _, t := range u.Todos {
		if t.Id == id {
			t.MarkDone()
			found = true
			return nil
		}
	}
	if !found {
		return errors.New("Can not find error")
	}
	return nil
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
