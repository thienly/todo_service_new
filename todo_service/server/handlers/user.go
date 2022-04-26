package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"new_todo_project/internal/database"
	"time"
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

type Todo struct {
	Title string `json:"title"`
	Content string `json:"content"`
	ExpiredDate time.Time `json:"expired_date"`
}
type UserRequest struct {
	Request
	Email string `json:"email"`
}
type UserResponse struct {
	Response
	Email    string `json:"email"`
	Name     string `json:"name"`
	Todos    []Todo `json:"todos"`
}
type AddTodoRequest struct {
	Request
	Email string `json:"email"`
	Todo Todo `json:"todo"`
}

type AddTotoResponse struct {
	Response
}
func (u *userHandlerImpl) Get(ctx context.Context, writer http.ResponseWriter, request *http.Request) {
	all, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	getRequest := &UserRequest{}
	err = json.Unmarshal(all, getRequest)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}

	panic("implement me")
}

func (u *userHandlerImpl) AddTodo(ctx context.Context, writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (u *userHandlerImpl) DeleteTodo(ctx context.Context, writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}
