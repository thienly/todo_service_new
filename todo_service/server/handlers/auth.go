package handlers

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"new_todo_project/internal/database"
	"new_todo_project/internal/domain"
	"time"
)

const (
	todoKey = "todo"
)

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type TokenResponse struct {
	Token string `json:"token"`
}
type AuthHandler interface {
	GenerateToken(writer http.ResponseWriter, request *http.Request)
	Register(writer http.ResponseWriter, request *http.Request)
}
type authHandlerImpl struct {
	client *mongo.Client
	TodoDb database.TodoDb
}

func NewAuthHandler(client *mongo.Client) AuthHandler {
	db := client.Database("todo")
	todo := database.NewTodoDb(db)
	return &authHandlerImpl{
		client: client,
		TodoDb: todo,
	}
}
func (u *authHandlerImpl) Register(writer http.ResponseWriter, request *http.Request) {
	all, err := io.ReadAll(request.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	userRes := &UserRegisterRequest{}
	err = json.Unmarshal(all, userRes)

	user := domain.NewUser(userRes.Name, userRes.Email, userRes.Password)
	ctx := context.Background()
	_, err = u.TodoDb.Add(ctx, user)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusCreated)
}

func (u *authHandlerImpl) GenerateToken(writer http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if err != nil {
		log.Err(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	tRequest := &TokenRequest{}
	err = json.Unmarshal(data, tRequest)
	if err != nil {
		log.Err(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// verify token request.

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(todoKey))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := TokenResponse{
		Token: tokenString,
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(responseData)
}
