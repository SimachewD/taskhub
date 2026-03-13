package service

import (
	"context"

	pb "github.com/SimachewD/taskhub/proto"
	"github.com/SimachewD/taskhub/internal/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
	pb.UnimplementedTaskServiceServer
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	return s.taskRepo.CreateTask(ctx, req)
}

func (s *TaskService) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
	return s.taskRepo.UpdateTask(ctx, req)
}

func (s *TaskService) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error) {
	return s.taskRepo.GetTask(ctx, req.Id)
}

func (s *TaskService) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Cursor == "" {
		req.Cursor = "infinity" //default pointer
	}

	tasks, nextCursor, err := s.taskRepo.ListTasks(ctx, req.UserId, req.Limit, req.Cursor)
	if err != nil {
		return nil, err
	}

	return &pb.ListTasksResponse{
		Tasks:      tasks,
		NextCursor: nextCursor,
	}, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	err := s.taskRepo.DeleteTask(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteTaskResponse{Success: true}, nil
}