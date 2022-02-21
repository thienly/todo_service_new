package main

import (
	"fmt"
	"net"
	"new_todo_project/pb"
	d "new_todo_project/pkg/database"
	"new_todo_project/pkg/services"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func startGrpcServer(logger zerolog.Logger, db *sqlx.DB, port int) error {
	todoDatabase := d.NewTodoDatabase(logger, db)
	todoService := services.NewTodoService(logger, todoDatabase)
	userService := services.NewUserService(logger, todoDatabase)
	loginService := services.NewLoginService(logger, todoDatabase)
	server := grpc.NewServer()
	pb.RegisterTodoServiceServer(server, todoService)
	pb.RegisterUserServiceServer(server, userService)
	pb.RegisterLoginServiceServer(server, loginService)
	lis, _ := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	logger.Info().Msg(fmt.Sprintf("Listening grpc on %s", fmt.Sprintf("0.0.0.0:%d", port)))
	err := server.Serve(lis)
	return err
}
