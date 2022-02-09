package domain_test

import (
	"github.com/stretchr/testify/require"
	"new_todo_project/pkg/domain"
	"testing"
)

func TestNewToDo(t *testing.T) {
	todo := domain.NewTodo("title", "description", false)
	if todo == nil {
		t.Fail()
	}
}
func TestMarkDone(t *testing.T) {
	todo := domain.NewTodo("Title", "Description", true)
	todo.MarkDone()
	if todo.Done == false {
		t.Fail()
	}

}

func TestNewUser(t *testing.T) {
	user := domain.NewUser("Test", "a@gmail.com", "a@gmail.com")
	if user == nil {
		t.Fail()
	}
}

func TestUser_GenerateToken(t *testing.T) {
	user := domain.NewUser("Test", "a@gmail.com", "a@gmail.com")
	require.Equal(t, user.GenerateToken(), "test_a@gmail.com", "should be equal")
}