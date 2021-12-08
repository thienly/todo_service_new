package services

import (
	"context"
	"github.com/rs/zerolog"
	"new_todo_project/internal/domain"
	"new_todo_project/pb"
	d "new_todo_project/pkg/database"
	"strconv"
)

type todoServiceImpl struct {
	logger zerolog.Logger
	db     d.TodoDatabase
	pb.UnimplementedTodoServiceServer
}

func NewTodoService(logger zerolog.Logger, db d.TodoDatabase) pb.TodoServiceServer {
	logger = logger.With().Str("package", "services").Logger()
	return &todoServiceImpl{
		db:     db,
		logger: logger,
	}
}

func (t *todoServiceImpl) Create(ctx context.Context, request *pb.TodoRequest) (*pb.TodoResponse, error) {
	todo:= domain.ToDomainTodo(request)
	id, err := t.db.AddNewTodo(ctx, int(request.UserId), todo)
	if err != nil {
		t.logger.Log().Err(err)
	}
	return &pb.TodoResponse{Id: strconv.Itoa(id)},nil
}


func (t *todoServiceImpl) MarkDone(ctx context.Context, request *pb.TodoMarkdoneRequest) (*pb.TodoMarkDoneResponse, error) {
	return nil,nil
}

func (t *todoServiceImpl) mustEmbedUnimplementedTodoServiceServer() {
}