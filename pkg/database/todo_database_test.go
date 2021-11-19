package todo_database

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestGetUsers(t *testing.T) {
	logger:= zerolog.New(os.Stdout)
	convey.Convey("the sql query should be called", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		cl := []string{"user_id", "username", "password", "email", "id","title","done"}
		rows := sqlmock.NewRows(cl)
		rows.AddRow(1, "test", "password", "p@password.com", nil, nil,nil)
		mock.ExpectQuery(GET_ALL_USER_QUERY).WillReturnRows(rows)
		database := NewTodoDatabase(logger, db)
		users, err := database.GetUsers(context.Background())
		if err != nil {
			t.Fail()
		}
		if users[0].Id != 1 {
			t.Fail()
		}
	})
}
