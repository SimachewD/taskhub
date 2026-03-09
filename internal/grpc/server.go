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
	pb.UnimplementedUserServiceServer
}

func NewGRPCServer(userRepo *repository.UserRepository) *GRPCServer {
	return &GRPCServer{userRepo: userRepo}
}

func (s *GRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	return s.userRepo.CreateUser(ctx, req)
}

func (s *GRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	return s.userRepo.GetUser(ctx, req.Id)
}

func StartServer(port string, userRepo *repository.UserRepository) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, NewGRPCServer(userRepo))

	log.Printf("gRPC server running on %s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}