package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type TokenResponse struct {
	Token   string    `json:"Token"`
	Expired time.Time `json:"expired_at"`
	Message string    `json:"message"`
}
type TokenRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

var (
	ErrInvalidToken = errors.New("Error while getting token. Please try again. ")
	todoService     = "https://todos.myshare.vn/"
)

//go:generate mockery --name TokenGetter --case=underscore
type TokenGetter interface {
	Get(request TokenRequest) (*TokenResponse, error)
}
type impl struct {
	httpClient *http.Client
}

func NewTokeGetter() TokenGetter {
	return &impl{
		httpClient: http.DefaultClient,
	}
}

func (tokenGetter *impl) Get(request TokenRequest) (*TokenResponse, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return nil, ErrInvalidToken
	}
	response, err := tokenGetter.httpClient.Post(todoService, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, ErrInvalidToken
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, ErrInvalidToken
	}
	data := &TokenResponse{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return data, nil
}
