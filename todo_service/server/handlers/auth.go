package handlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

const (
	todoKey = "todo"
)

type TokenRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
type TokenResponse struct {
	Token   string    `json:"token"`
}
func GenerateToken(writer http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if err != nil {
		log.Err(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	trequest:= &TokenRequest{}
	err = json.Unmarshal(data, trequest)
	if err != nil {
		log.Err(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(todoKey))
	if err != nil {
		writer.Write([]byte(err.Error()))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := TokenResponse{
		Token:   tokenString,
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(responseData)
}