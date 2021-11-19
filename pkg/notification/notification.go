package notification

import (
	"new_todo_project/pkg/domain"
)

type Notifier interface {
	Notify(u *domain.User, todo *domain.Todo)
}

type notifierImpl struct {
}

// Notify Sending new email to end-user if there is an overdue todo list.
func (notification *notifierImpl) Notify(u *domain.User, todo *domain.Todo) {

}
