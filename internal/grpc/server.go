package grpc

import (
	"context"
	"log"
	"net"

	pb "github.com/SimachewD/taskhub/proto"
	"github.com/SimachewD/taskhub/internal/repository"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	userRepo *repository.UserRepository
	taskRepo *repository.TaskRepository
	pb.UnimplementedUserServiceServer
	pb.UnimplementedTaskServiceServer
}

func NewGRPCServer(userRepo *repository.UserRepository, taskRepo *repository.TaskRepository) *GRPCServer {
	return &GRPCServer{userRepo: userRepo, taskRepo: taskRepo}
}

// user services
func (s *GRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	return s.userRepo.CreateUser(ctx, req)
}

func (s *GRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	return s.userRepo.GetUser(ctx, req.Id)
}

// task services
func (s *GRPCServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	return s.taskRepo.CreateTask(ctx, req)
}

func (s *GRPCServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error) {
	return s.taskRepo.GetTask(ctx, req.Id)
}

func (s *GRPCServer) ListTasks(
	ctx context.Context,
	req *pb.ListTasksRequest,
) (*pb.ListTasksResponse, error) {

	if req.Limit == 0 {
		req.Limit = 10
	}

	tasks, nextCursor, err := s.taskRepo.ListTasks(
		ctx,
		req.UserId,
		req.Limit,
		req.Cursor,
	)

	if err != nil {
		return nil, err
	}

	return &pb.ListTasksResponse{
		Tasks:      tasks,
		NextCursor: nextCursor,
	}, nil
}

func StartServer(port string, userRepo *repository.UserRepository, taskRepo *repository.TaskRepository) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	grpcServer := NewGRPCServer(userRepo, taskRepo)

	pb.RegisterUserServiceServer(server, grpcServer)
	pb.RegisterTaskServiceServer(server, grpcServer)

	log.Printf("gRPC server running on %s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}