package service

import (
	"context"

	pb "github.com/SimachewD/taskhub/proto"
	"github.com/SimachewD/taskhub/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
	pb.UnimplementedUserServiceServer
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	return s.userRepo.CreateUser(ctx, req)
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	return s.userRepo.GetUser(ctx, req.Id)
}