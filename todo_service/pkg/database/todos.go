package todo_database

import (
	"context"
	"new_todo_project/internal/domain"
)

func (t *todoImpl) AddNewTodo(ctx context.Context, userId int, todo *domain.Todo) (int, error) {
	rows:= t.db.QueryRowContext(ctx, INSERT_NEW_TODO,userId, todo.Title,todo.Description, todo.Done)
	var id int
	rows.Scan(&id)
	return id, nil
}

func (t *todoImpl) MarkDone(ctx context.Context, u *domain.User, id int) error {
	panic("implement me")
}
