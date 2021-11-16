package todo_database

import (
	"context"
	"database/sql"
	"new_todo_project/pkg/domain"
)

type TodoDatabase interface {
	AddNewUser(ctx context.Context, u *domain.User) (int, error)
	AddNewTodo(ctx context.Context, u *domain.User, todo *domain.Todo) error
	MarkDone(ctx context.Context, u *domain.User, todoId int) error
}

type todoImpl struct {
	db *sql.DB
}

func NewTodoDatabase(pool *sql.DB) TodoDatabase {
	return &todoImpl{
		db: pool,
	}
}

func (t *todoImpl) AddNewUser(ctx context.Context, u *domain.User) (int, error) {
	panic("implement me")
}

func (t *todoImpl) AddNewTodo(ctx context.Context, u *domain.User, todo *domain.Todo) error {
	panic("implement me")
}

func (t *todoImpl) MarkDone(ctx context.Context, u *domain.User, id int) error {
	panic("implement me")
}
