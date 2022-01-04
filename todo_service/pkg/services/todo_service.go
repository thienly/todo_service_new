package services

import (
	"context"
	"github.com/rs/zerolog"
	"new_todo_project/internal/domain"
	"new_todo_project/pb"
	"strconv"

	db "new_todo_project/pkg/database"
)

type todoServiceImpl struct {
	logger zerolog.Logger
	db     db.TodoDatabase
	pb.UnimplementedTodoServiceServer
}

func NewTodoService(logger zerolog.Logger, db db.TodoDatabase) pb.TodoServiceServer {
	return &todoServiceImpl{
		db: db,
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
	err := t.db.MarkDone(ctx, int(request.TodoId))
	if err != nil {
		return nil, err
	}
	return nil, err
}

