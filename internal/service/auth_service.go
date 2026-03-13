package service

import (
	"context"
	"github.com/SimachewD/taskhub/internal/repository"
	pb "github.com/SimachewD/taskhub/proto"
)

type AuthService struct {
	userRepo *repository.UserRepository
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	return s.userRepo.Register(ctx, req.Name, req.Email, req.Password)
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	return s.userRepo.Login(ctx, req.Email, req.Password)
}