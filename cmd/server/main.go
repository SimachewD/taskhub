package main

import (
	"context"
	"log"

	"github.com/SimachewD/taskhub/internal/cache"
	"github.com/SimachewD/taskhub/internal/config"
	"github.com/SimachewD/taskhub/internal/database"
	"github.com/SimachewD/taskhub/internal/grpc"
	"github.com/SimachewD/taskhub/internal/repository"
	"github.com/SimachewD/taskhub/internal/worker"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}
	
	ctx := context.Background()

	// 1. Load config (contains DB DSN)
    cfg := config.Load()

    // 2. Run migrations
    err := database.RunMigrations(cfg.PostgresURL)
    if err != nil {
        log.Fatalf("failed to run migrations: %v", err)
    }

    log.Println("migrations ran successfully")

	dbPool := database.NewPostgresPool(ctx, cfg.PostgresURL)
	redisClient := cache.NewRedisClient(cfg.RedisAddr)

	userRepo := repository.NewUserRepository(dbPool, redisClient)
	taskRepo := repository.NewTaskRepository(dbPool, redisClient)

	worker := worker.NewWorker(redisClient)

	go worker.Start(ctx)

	grpc.StartServer(":50051", userRepo, taskRepo, redisClient)
}