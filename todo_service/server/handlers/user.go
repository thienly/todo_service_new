package handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"new_todo_project/internal/database"
)

type UserHandler interface {
	Get(ctx context.Context, writer http.ResponseWriter, request *http.Request)
	AddTodo(ctx context.Context, writer http.ResponseWriter, request *http.Request)
	DeleteTodo(ctx context.Context, writer http.ResponseWriter, request *http.Request)
}

type userHandlerImpl struct {
	client *mongo.Client
	TodoDb database.TodoDb
}

func NewUserHandler(client *mongo.Client) UserHandler {
	db := client.Database("todo")
	todo := database.NewTodoDb(db)
	return &userHandlerImpl{
		client: client,
		TodoDb: todo,
	}
}
func (u *userHandlerImpl) Get(ctx context.Context, writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (u *userHandlerImpl) AddTodo(ctx context.Context, writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (u *userHandlerImpl) DeleteTodo(ctx context.Context, writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}
