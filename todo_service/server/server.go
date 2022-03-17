package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"new_todo_project/server/config"
	"new_todo_project/server/handlers"
	"os"
	"os/signal"
)

type Server struct {
	AppConfig *config.Config
	ReadyChan chan bool
}
func (s *Server) Start() {
	clientOpts := options.Client().ApplyURI(s.AppConfig.Db.Conn)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		panic(fmt.Sprintf("Can not connect to db:%v", err))
	}
	r:= mux.NewRouter()
	authHandler:= handlers.NewAuthHandler(client)
	s.Handler(r, authHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", r))
	}()
	if s.ReadyChan != nil {
		s.ReadyChan <- true
	}
	log.Println("Application is started and listening on 8080")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// wait to consume signal
	<-c
	log.Println("Application is shutting down")
	os.Exit(1)
}

func (s *Server) Handler(router *mux.Router, authHandler handlers.AuthHandler) {
	router.HandleFunc("/token", authHandler.GenerateToken).Methods("POST")
	router.HandleFunc("/users/register", authHandler.Register).Methods("POST")
}
