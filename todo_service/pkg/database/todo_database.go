package todo_database

import (
	"context"
	"database/sql"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/rs/zerolog"
	"new_todo_project/internal/domain"
)

type DatabaseResult struct {
	data interface{}
	err  error
}

const (
	GET_ALL_USER_QUERY = `SELECT u.user_id, u.username, u.password, u.email, t.id, t.title, t.done
						  FROM public.users u
						  LEFT JOIN todos t on u.user_id = t.user_id`
	INSERT_NEW_USER = "INSERT INTO public.users(username, password, email) VALUES ($1,$2,$3) returning user_id"
	INSERT_NEW_TODO = `INSERT INTO public.todos(user_id,title, description, done) VALUES ($1, $2, $3, $4) returning id`
)

type TodoDatabase interface {
	AddNewUser(ctx context.Context, u *domain.User) (int64, error)
	AddNewTodo(ctx context.Context, userId int, todo *domain.Todo) (int, error)
	MarkDone(ctx context.Context, u *domain.User, todoId int) error
	GetUsers(ctx context.Context) ([]*domain.User, error)
}

type todoImpl struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewTodoDatabase(logger zerolog.Logger, pool *sql.DB) TodoDatabase {
	logger = logger.With().Str("package", "todo_database").Logger()
	hystrix.ConfigureCommand("database", hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests: 100,
		SleepWindow:           100,
		ErrorPercentThreshold: 20,
	})
	return &todoImplProxy{
		todoImpl{
			db:     pool,
			logger: logger,
		},
	}
}

type todoImplProxy struct {
	todoImpl
}

func (proxy *todoImplProxy) AddNewUser(ctx context.Context, u *domain.User) (int64, error) {
	resultChan := make(chan int64)
	errChan := hystrix.Go("database", func() error {
		user, err := proxy.todoImpl.AddNewUser(ctx, u)
		resultChan <- user
		return err
	}, func(err error) error {
		if ok := retryAbleError(err); ok {
			_, _ = proxy.AddNewUser(ctx, u)
		}
		return err
	})
	select {
	case out := <-resultChan:
		return out, nil
	case err := <-errChan:
		return 0, err
	}
}

func retryAbleError(err error) bool {
	return true
}

func (proxy *todoImplProxy) AddNewTodo(ctx context.Context, userId int, todo *domain.Todo) (int, error) {
	resultChan := make(chan int)
	errChan := hystrix.Go("database", func() error {
		user, err := proxy.todoImpl.AddNewTodo(ctx, userId, todo)
		resultChan <- user
		return err
	}, func(err error) error {
		if ok := retryAbleError(err); ok {
			_, _ = proxy.AddNewTodo(ctx, userId, todo)
		}
		return err
	})
	select {
	case out := <-resultChan:
		return out, nil
	case err := <-errChan:
		return 0, err
	}

	HytrixWrapper(resultChan,errChan, func() (interface{}, error) {
		newTodo, err := proxy.todoImpl.AddNewTodo(ctx, userId, todo)
		return newTodo, err
	})
}
// delay call
// high order function
func abc(run func() error) func(i interface{}){
	return func(i interface{}) {
		run()
	}
}
// high order function

func HytrixWrapper(result chan interface{}, errChan chan error, run func() (interface{},error)) {
	errChan = hystrix.Go("database", func() error {
		// logging
		// tracing
		i, err := run()
		result <- i
		return err
	}, func(err error) error {
		_, _ = run()
		return err
	})
}