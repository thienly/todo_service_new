package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"new_todo_project/pb"
)

func startHttpServer(log zerolog.Logger, grpcPort, port int) error{
	ctx := context.Background()
	ctxCancel, _ := context.WithCancel(ctx)
	gwMux := runtime.NewServeMux()
	opts := grpc.WithInsecure()
	err := pb.RegisterTodoServiceHandlerFromEndpoint(ctxCancel, gwMux, fmt.Sprintf("0.0.0.0:%d", grpcPort), []grpc.DialOption{opts})
	err = pb.RegisterUserServiceHandlerFromEndpoint(ctxCancel, gwMux,fmt.Sprintf("0.0.0.0:%d", grpcPort), []grpc.DialOption{opts})
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
	log.Info().Msg(fmt.Sprintf("Listening httpServer on %s", fmt.Sprintf("0.0.0.0:%d", port)))
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), swaggerMux)
}
