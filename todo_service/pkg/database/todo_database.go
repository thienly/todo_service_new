package todo_database

import (
	"context"
	"database/sql"
	"new_todo_project/internal/domain"

	"github.com/rs/zerolog"
)

type DatabaseResult struct {
	data interface{}
	err  error
}

const (
	GET_ALL_USER_QUERY = `SELECT u.id, u.username, u.password, u.email, t.id, t.title, t.done
						  FROM public.users u
						  LEFT JOIN todos t on u.id = t.user_id`
	INSERT_NEW_USER = "INSERT INTO public.users(username, password, email) VALUES ($1,$2,$3) returning user_id"
	INSERT_NEW_TODO = `INSERT INTO public.todos(user_id,title, description, done) VALUES ($1, $2, $3, $4) returning id`
	MARK_DONE_TODO = `UPDATE public.todos SET done = true WHERE id = $1`
)

type TodoDatabase interface {
	AddNewUser(ctx context.Context, u *domain.User) (int64, error)
	AddNewTodo(ctx context.Context, userId int, todo *domain.Todo) (int, error)
	MarkDone(ctx context.Context, todoId int) error
	GetUsers(ctx context.Context) ([]*domain.User, error)
}

type todoImpl struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewTodoDatabase(logger zerolog.Logger, pool *sql.DB) TodoDatabase {
	logger = logger.With().Str("package", "todo_database").Logger()
	return &todoImpl{
		db:     pool,
		logger: logger,
	}
}
