package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"new_todo_project/server/handlers"
	"os"
	"os/signal"
)

func Start() {
	log.Println("Starting application!!")
	r := mux.NewRouter()

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
func Handler(router mux.Router) {
	router.HandleFunc("/token", handlers.GenerateToken).Methods("POST")
	router.HandleFunc("/users/register", handlers.Register).Methods("POST")
}
