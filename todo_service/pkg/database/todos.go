package todo_database

import (
	"context"
	"new_todo_project/internal/domain"
)

func (t *todoImpl) AddNewTodo(ctx context.Context, userId int, todo *domain.Todo) (int, error) {
	rows:= t.db.QueryRowContext(ctx, INSERT_NEW_TODO,userId, todo.Title,todo.Description, todo.Done)
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
