package services

import (
	"context"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"new_todo_project/internal/domain"
	"new_todo_project/pb"
	db "new_todo_project/pkg/database"
)

type userServiceImpl struct {
	logger zerolog.Logger
	db     db.TodoDatabase
	pb.UnimplementedUserServiceServer
}

func NewUserService(logger zerolog.Logger, db db.TodoDatabase) pb.UserServiceServer {
	return &userServiceImpl{
		db:     db,
		logger: logger,
	}
}

func (u *userServiceImpl) CreateUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	u.logger.Info().Msg("Start create user")
	user := &domain.User{
		Name:     in.GetName(),
		Password: in.GetPassword(),
		Email:    in.GetEmail(),
	}
	dbResult, err := u.db.AddNewUser(ctx, user)
	if err != nil {
		u.logger.Err(err)
		return nil, err
	}
	u.logger.Info().Msg("Ended create users")
	return &pb.UserResponse{UserId: dbResult}, nil
}

func (u *userServiceImpl) GetAllUsers(ctx context.Context, empty *emptypb.Empty) (*pb.UserList, error) {
	u.logger.Info().Msg("Get users")
	users, err := u.db.GetUsers(ctx)
	if err != nil {
		u.logger.Log().Err(err)
		return nil, status.Errorf(codes.Internal, "Internal server error!")
	}
	var result pb.UserList
	for _, user := range users {
		result.Users = append(result.Users, domain.ToPbUser(user))
	}
	u.logger.Info().Msg("End Get users")
	return &result, nil
}
func (t *userServiceImpl) mustEmbedUnimplementedTodoServiceServer() {
}
