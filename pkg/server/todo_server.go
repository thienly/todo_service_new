package server

import (
	"context"
	"new_todo_project/pb"

)

type TodoServer struct {
	pb.UnimplementedTodoServiceServer
}

func (t TodoServer) Create(ctx context.Context, request *pb.TodoRequest) (*pb.TodoResponse, error) {
	panic("implement me")
}

func (t TodoServer) GetAll(ctx context.Context, void *pb.Void) (*pb.TodoList, error) {
	panic("implement me")
}

func (t TodoServer) MarkDone(ctx context.Context, request *pb.TodoRequest) (*pb.TodoResponse, error) {
	panic("implement me")
}

func (t TodoServer) mustEmbedUnimplementedTodoServiceServer() {
	panic("implement me")
}

func NewTodoServer() pb.TodoServiceServer {
	return &TodoServer{}
}
