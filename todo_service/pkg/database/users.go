package todo_database

import (
	"context"
	"database/sql"
	"errors"
	pg "github.com/lib/pq"
	"new_todo_project/internal/domain"
)
// GetUsers gets all users of system
func (t *todoImpl) GetUsers(ctx context.Context) ([]*domain.User, error) {
	rows, err := t.db.QueryContext(ctx, GET_ALL_USER_QUERY)
	defer func() {
		if r := recover(); r != nil {
			t.logger.Err(err)
		}
		_: rows.Close()
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
