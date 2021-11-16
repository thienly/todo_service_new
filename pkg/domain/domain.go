package domain

import "errors"

type Todo struct {
	Id          int
	Title       string
	Description string
	Done        bool
}

func NewTodo(title, description string, done bool) Todo {
	return Todo{
		Title:       title,
		Description: description,
		Done:        done,
	}
}

func (t *Todo) MarkDone() {
	t.Done = true
}

type User struct {
	Id       int
	Name     string
	Password string
	Email    string
}

func NewUser(name, password, email string) User {
	return User{Name: name,
		Password: password,
		Email:    email}
}

func (u *User) ChangePassword(oldPassword, newPassword string) (bool, error){
	if u.Password != oldPassword{
		return false, errors.New("Old password is not matched with new password")
	}
	u.Password = newPassword
	return true, nil
}
