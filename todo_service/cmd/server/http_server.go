package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"new_todo_project/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startHttpServer(log zerolog.Logger, grpcPort, port int) error {
	serverAddr := fmt.Sprintf("0.0.0.0:%d", grpcPort)
	ctx := context.Background()
	gwMux := runtime.NewServeMux()
	if err:= registerEndpoints(ctx, serverAddr, gwMux); err != nil {
		return err
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
func registerEndpoints(ctx context.Context, serverAddr string, gwMux *runtime.ServeMux) error {
	ctxCancel, _ := context.WithCancel(ctx)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterTodoServiceHandlerFromEndpoint(ctxCancel, gwMux, serverAddr, opts)
	err = pb.RegisterUserServiceHandlerFromEndpoint(ctxCancel, gwMux, serverAddr, opts)
	err = pb.RegisterLoginServiceHandlerFromEndpoint(ctxCancel, gwMux, serverAddr, opts )
	return err
}