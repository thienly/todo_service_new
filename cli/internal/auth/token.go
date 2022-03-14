package auth

import "time"

type TokenResponse struct {
	Data string `json:"data"`
	Expired time.Time `json:"expired_at"`
}
//go:generate mockery --name TokenGetter --case=underscore
type TokenGetter interface {
	Get(username, password string) (TokenResponse, error)
}
