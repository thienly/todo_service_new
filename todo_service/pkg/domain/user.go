package domain

import (
	"errors"
	"fmt"
	"strings"
)

type User struct {
	Id       int    `sql-col:"id"`
	Name     string `sql-col:"username"`
	Password string `sql-col:"password"`
	Email    string `sql-col:"email"`
	Todos    []*Todo
}

func NewUser(name, password, email string) *User {
	return &User{
		Name: name,
		Password: password,
		Email:    email,
	}
}

func (u *User) ChangePassword(oldPassword, newPassword string) (bool, error) {
	if u.Password != oldPassword {
		return false, errors.New("old password is not matched with new password")
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

func (u *User) GenerateToken() string{
	token:= fmt.Sprintf("%v_%v", strings.ToLower(u.Name), strings.ToLower(u.Password))
	return token
}