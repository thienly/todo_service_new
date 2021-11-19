package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"net/http"
	"new_todo_project/pb"
	"new_todo_project/pkg/config"
	d "new_todo_project/pkg/database"
	"new_todo_project/pkg/services"
	"os"
	"os/signal"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Starting application!!")
	lis, _ := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 25000))
	_, err := config.LoadFromJsonOrPanic("config.json")
	if err != nil {
		logger.Fatal().Msg("Can not read config.json file!!")
	}
	db, err := sql.Open("postgres", config.AppCfg.Database.ConnectionString)
	err = db.Ping()
	if err != nil {
		logger.Fatal().Msg("Can not open database connection")
	}
	defer db.Close()
	todoDatabase := d.NewTodoDatabase(logger, db)
	todoService := services.NewTodoService(logger, todoDatabase)
	userService := services.NewUserService(logger, todoDatabase)

	go func() {
		server := grpc.NewServer()
		pb.RegisterTodoServiceServer(server, todoService)
		pb.RegisterUserServiceServer(server, userService)
		err := server.Serve(lis)
		if err != nil {
			logger.Fatal().Err(err)
			return
		}
	}()

	go func() {
		startHttpServer()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Log().Msg("Shutdown application")
	os.Exit(0)
}

func startHttpServer() {
	ctx := context.Background()
	ctxCancel, _ := context.WithCancel(ctx)
	gwMux := runtime.NewServeMux()
	opts := grpc.WithInsecure()
	err := pb.RegisterTodoServiceHandlerFromEndpoint(ctxCancel, gwMux, "0.0.0.0:25000", []grpc.DialOption{opts})
	err = pb.RegisterUserServiceHandlerFromEndpoint(ctxCancel, gwMux, "0.0.0.0:25000", []grpc.DialOption{opts})
	if err != nil {
		log.Fatal().Msg("Can not establish connection to grpc server")
	}
	// swagger server.
	swaggerMux := http.NewServeMux()
	swaggerMux.Handle("/", gwMux)
	swaggerMux.HandleFunc("/swagger.json", func(writer http.ResponseWriter, request *http.Request) {
		data, _ := ioutil.ReadFile("swagger/v1/todo.swagger.json")
		writer.Write(data)
	})
	fs := http.FileServer(http.Dir("cmd/server/www/swagger-ui"))
	swaggerMux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", fs))
	http.ListenAndServe("0.0.0.0:8081", swaggerMux)
}
