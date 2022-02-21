package todo_database

import (
	"context"
	"database/sql"
	"errors"
	"new_todo_project/pkg/domain"

	"github.com/jmoiron/sqlx"
	pg "github.com/lib/pq"
	"github.com/rs/zerolog"
)

const (
	GET_ALL_USER_QUERY = `SELECT u.id, u.username, u.password, u.email, t.id, t.title, t.done
						  FROM public.users u
						  LEFT JOIN todos t on u.id = t.user_id`
	INSERT_NEW_USER           = "INSERT INTO public.users(username, password, email) VALUES ($1,$2,$3) returning user_id"
	INSERT_NEW_TODO           = `INSERT INTO public.todos(user_id,title, description, done) VALUES ($1, $2, $3, $4) returning id`
	MARK_DONE_TODO            = `UPDATE public.todos SET done = true WHERE id = $1`
	GET_USER_BY_USERNAME_PASS = `SELECT u.id, u.username, u.password from public.users u where u.username = $1 and u.password = $2`
)

type TodoDatabase interface {
	AddNewUser(ctx context.Context, u *domain.User) (int64, error)
	AddNewTodo(ctx context.Context, userId int, todo *domain.Todo) (int, error)
	MarkDone(ctx context.Context, todoId int) error
	GetUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByUserNameAndPassword(ctx context.Context, userName string, password string) (*domain.User, error)
}

type todoImpl struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

func NewTodoDatabase(logger zerolog.Logger, pool *sqlx.DB) TodoDatabase {
	logger = logger.With().Str("package", "todo_database").Logger()
	return &todoImpl{
		db:     pool,
		logger: logger,
	}
}

func (t *todoImpl) AddNewTodo(ctx context.Context, userId int, todo *domain.Todo) (int, error) {
	rows := t.db.QueryRowContext(ctx, INSERT_NEW_TODO, userId, todo.Title, todo.Description, todo.Done)
	var id int
	err := rows.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (t *todoImpl) MarkDone(ctx context.Context, id int) error {
	t.logger.Info().Int("Id", id).Msg("todo mark done")
	_, err := t.db.ExecContext(ctx, MARK_DONE_TODO, id)
	return err
}

// GetUsers gets all users of system.
func (t *todoImpl) GetUsers(ctx context.Context) ([]*domain.User, error) {

	rows, err := t.db.QueryContext(ctx, GET_ALL_USER_QUERY)
	defer func() {
		if r := recover(); r != nil {
			t.logger.Err(err)
		}
	}()
	if err != nil {
		pgErr := err.(*pg.Error)
		t.logger.Error().Err(pgErr).Msg(pgErr.Message)
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

func (t *todoImpl) AddNewUser(ctx context.Context, u *domain.User) (int64, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err == nil {
			err := tx.Commit()
			if err != nil {
				t.logger.Log().Err(err).Msg("Can not commit transaction. transaction aborted.")
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
				err = errors.New("duplication data")
				return 0, err
			} else {
				err = errors.New("internal server")
				return 0, err
			}
		} else {
			return 0, rows.Err()
		}
	}
	err = rows.Scan(&id)
	return id, nil
}

func (c *todoImpl) GetUserByUserNameAndPassword(ctx context.Context, userName, password string) (*domain.User, error) {
	rows, err := c.db.QueryContext(ctx, GET_USER_BY_USERNAME_PASS, userName, password)
	if err != nil {
		return nil, errors.New("can not find user")
	}

	rows.Next()
	var id sql.NullInt64
	var usernameCl sql.NullString
	var passwordCl sql.NullString
	err = rows.Scan(&id, &usernameCl, &passwordCl)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Id:       int(id.Int64),
		Name:     usernameCl.String,
		Password: passwordCl.String,
	}
	return user, nil
}
