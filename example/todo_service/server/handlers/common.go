package handlers

import (
	"net/http"
)

type Request struct {
	RequestId string `json:"request_id"`
	Endpoint string `json:"endpoint"`
}

type ErrorMessage struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	RequestId string `json:"request_id"`
	ErrorMessages []ErrorMessage `json:"error_messages"`
}


func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log before receive request
		next.ServeHTTP(w, r)
		// log after process request.
	})

}