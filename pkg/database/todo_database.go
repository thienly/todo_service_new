package todo_database

import (
	"context"
	"database/sql"
	"errors"
	pg "github.com/lib/pq"
	"github.com/rs/zerolog"
	"new_todo_project/pkg/domain"
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
	AddNewTodo(ctx context.Context, userId int, todo *domain.Todo) (int,error)
	MarkDone(ctx context.Context, u *domain.User, todoId int) error
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

func (t *todoImpl) AddNewUser(ctx context.Context, u *domain.User) (int64, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err == nil {
			err := tx.Commit()
			if err != nil {
				t.logger.Log().Err(err).Msg("Can not commit transaction. transation aborted.")
			}
		} else {
			err := tx.Rollback()
			if err != nil {
				t.logger.Log().Err(err).Msg("Can not rollback transaction")
			}
		}
	}()
	var id int64
	rows := tx.QueryRowContext(ctx, INSERT_NEW_USER, u.Name, u.Password, u.Email)
	if rows.Err() != nil {
		rawError, ok := rows.Err().(*pg.Error)
		if ok {
			if rawError.Code == ("23505") {
				err = errors.New("Duplication data")
				return 0, err
			} else {
				err = errors.New("Internal server")
				return 0, err
			}
		} else {
			return 0, rows.Err()
		}
	}
	rows.Scan(&id)
	return id, nil
}

func (t *todoImpl) AddNewTodo(ctx context.Context, userId int, todo *domain.Todo) (int, error) {
	rows:= t.db.QueryRowContext(ctx, INSERT_NEW_TODO,userId, todo.Title,todo.Description, todo.Done)
	var id int
	rows.Scan(&id)
	return id, nil
}

func (t *todoImpl) MarkDone(ctx context.Context, u *domain.User, id int) error {
	panic("implement me")
}

func (t *todoImpl) GetUsers(ctx context.Context) ([]*domain.User, error) {
	rows, err := t.db.QueryContext(ctx, GET_ALL_USER_QUERY)
	defer func() {
		if r := recover(); r != nil {
			t.logger.Err(err)
		}
		rows.Close()
	}()

	if err != nil {
		return nil, err
	}
	results := make(map[int]*domain.User)
	for rows.Next() {
		var u domain.User
		var todoId sql.NullInt64
		var todoTitle sql.NullString
		var todoDone sql.NullBool

		if err := rows.Scan(&u.Id, &u.Name, &u.Password, &u.Email, &todoId, &todoTitle, &todoDone); err != nil {
			t.logger.Err(err)
		}
		if results[u.Id] == nil {
			results[u.Id] = &u
		} else {
			if todoId.Valid {
				var t domain.Todo
				todoId.Scan(&t.Id)
				todoTitle.Scan(&t.Title)
				todoDone.Scan(&t.Done)
				results[u.Id].Todos = append(results[u.Id].Todos, &t)
			}
		}
	}
	var data []*domain.User
	for _, user := range results {
		data = append(data, user)
	}
	return data, nil

}
