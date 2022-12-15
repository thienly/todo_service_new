package domain

import "time"

type Todo struct {
	Id          string
	Title       string
	ExpiredData time.Time
	IsDone      bool
	CreatedDate time.Time
	CreatedBy   string
}

func (todo *Todo) MarkDone() {
	todo.IsDone = true
}
func NewTodo(title string, d time.Time, createdBy string) *Todo {
	return nil
}
