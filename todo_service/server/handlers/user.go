package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

type UserRegister struct {
	Name     string
	Email    string
	Password string
}

func Register(writer http.ResponseWriter, request *http.Request) {
	all, err := io.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	userRes := &UserRegister{}
	err = json.Unmarshal(all, userRes)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

}
