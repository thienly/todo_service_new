package main

import (
	"fmt"
	"new_todo_project/server"
	"new_todo_project/server/config"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("panic application: %v", err))
	}
	loadConfig, err := config.LoadConfig(path)
	if err != nil {
		panic(fmt.Sprintf("panic application: %v", err))
	}
	newServer := &server.Server{AppConfig: loadConfig}
	newServer.Start()
}
