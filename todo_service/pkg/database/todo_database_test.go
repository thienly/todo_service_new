package todo_database_test

import (
	"context"
	"errors"
	todo_database "new_todo_project/pkg/database"
	"new_todo_project/pkg/domain"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetUsers(t *testing.T) {
	logger := zerolog.New(os.Stdout)
	convey.Convey("the sql query should be called", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		cl := []string{"user_id", "username", "password", "email", "id", "title", "done"}
		rows := sqlmock.NewRows(cl)
		rows.AddRow(1, "test", "password", "p@password.com", nil, nil, nil)
		mock.ExpectQuery(todo_database.GET_ALL_USER_QUERY).WillReturnRows(rows)
		database := todo_database.NewTodoDatabase(logger, db)
		users, err := database.GetUsers(context.Background())
		if err != nil {
			t.Fail()
		}
		if users[0].Id != 1 {
			t.Fail()
		}
	})
}

func TestAddNewUser(t *testing.T) {
	convey.Convey("the query should be called", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fail()
		}
		_ = mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"user_id"})
		rows.AddRow("1")
		mock.ExpectQuery(todo_database.INSERT_NEW_USER).WithArgs("test", "test", "test").WillReturnRows().WillReturnRows(rows)
		mock.ExpectCommit()
		database := todo_database.NewTodoDatabase(zerolog.New(os.Stdout), db)
		id, err := database.AddNewUser(context.Background(), &domain.User{
			Name:     "test",
			Password: "test",
			Email:    "test",
		})
		if id != 1 {
			t.Fail()
		}
	})

	convey.Convey("should call rollback when there is an error", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fail()
		}
		mock.ExpectBegin()
		mock.ExpectRollback()
		mock.ExpectQuery(todo_database.INSERT_NEW_USER).WillReturnError(errors.New("There is an error"))
		database := todo_database.NewTodoDatabase(zerolog.New(os.Stdout), db)
		_, _ = database.AddNewUser(context.Background(), &domain.User{})
	})
}

func TestAddNewTodo(t *testing.T) {
	convey.Convey("should call sql", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fail()
		}
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"id"})
		rows.AddRow("1")
		mock.ExpectQuery(todo_database.INSERT_NEW_TODO).WillReturnRows(rows)
		todoDatabase := todo_database.NewTodoDatabase(zerolog.New(os.Stdout), db)
		_, _ = todoDatabase.AddNewTodo(context.Background(), 1, &domain.Todo{
			Title:       "",
			Description: "",
			Done:        false,
		})
	})
}

func TestMarkDone(t *testing.T) {
	convey.Convey("should call sql", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fail()
		}
		mock.ExpectExec(todo_database.MARK_DONE_TODO)
		database := todo_database.NewTodoDatabase(zerolog.New(os.Stdout), db)
		_ = database.MarkDone(context.Background(), 1)
	})
}
