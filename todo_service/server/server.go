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

func Start() {
	log.Println("Starting application!!")
	r := mux.NewRouter()
	loadConfig, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("The configuration is load loaded %v", err))
	}
	clientOpts := options.Client().ApplyURI(loadConfig.Db.Conn)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		panic(fmt.Sprintf("Can not connect to db:%v", err))
	}
	authHandler:= handlers.NewAuthHandler(client)
	Handler(r, authHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", r))
	}()
	log.Println("Application is started and listening on 8080")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Application is shutting down")
	os.Exit(1)
}
func Handler(router *mux.Router, authHandler handlers.AuthHandler) {
	router.HandleFunc("/token", authHandler.GenerateToken).Methods("POST")
	router.HandleFunc("/users/register", handlers.Register).Methods("POST")
}
