package grpc

import (
	"log"
	"net"

	"github.com/SimachewD/taskhub/internal/cache"
	"github.com/SimachewD/taskhub/internal/middleware"
	"github.com/SimachewD/taskhub/internal/repository"
	"github.com/SimachewD/taskhub/internal/service"
	pb "github.com/SimachewD/taskhub/proto"
	"google.golang.org/grpc"
)

func StartServer(port string, userRepo *repository.UserRepository, taskRepo *repository.TaskRepository, redis *cache.RedisClient) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(middleware.AuthInterceptor()),)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)
	authService := service.NewAuthService(userRepo)
	notificationService := service.NewNotificationServer(redis)

	// Register services
	pb.RegisterUserServiceServer(server, userService)
	pb.RegisterTaskServiceServer(server, taskService)
	pb.RegisterAuthServiceServer(server, authService)
	pb.RegisterNotificationServiceServer(server, notificationService)

	log.Printf("gRPC server running on %s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}