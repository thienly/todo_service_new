package services

import (
	"context"
	"github.com/rs/zerolog"
	"new_todo_project/pb"
	db "new_todo_project/pkg/database"
)

type LoginService struct {
	logger zerolog.Logger
	db     db.TodoDatabase
	pb.UnimplementedLoginServiceServer
}

func NewLoginService(logger zerolog.Logger, todo db.TodoDatabase) *LoginService {
	return &LoginService{
		logger: logger,
		db:     todo,
	}
}
func (login *LoginService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := login.db.GetUserByUserNameAndPassword(ctx, request.Username, request.Password)
	if err != nil {
		return nil, err
	}
	token := user.GenerateToken()
	return &pb.LoginResponse{Token: token}, nil
}
