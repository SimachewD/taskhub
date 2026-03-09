package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/SimachewD/taskhub/internal/cache"
	pb "github.com/SimachewD/taskhub/proto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db    *pgxpool.Pool
	redis *cache.RedisClient
}

func NewUserRepository(db *pgxpool.Pool, redis *cache.RedisClient) *UserRepository {
	return &UserRepository{db: db, redis: redis}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *pb.CreateUserRequest) (*pb.UserResponse, error) {
	id := uuid.NewString()
	_, err := r.db.Exec(ctx,
		"INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
		id, user.Name, user.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	resp := &pb.UserResponse{
		Id:    id,
		Name:  user.Name,
		Email: user.Email,
	}

	// Cache the new user for fast reads
	data, _ := json.Marshal(resp)
	_ = r.redis.Set(ctx, "user:"+id, data, 5*time.Minute)

	return resp, nil
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (*pb.UserResponse, error) {
	// Try Redis first
	if r.redis != nil {
		cached, err := r.redis.Get(ctx, "user:"+id)
		if err == nil {
			log.Println("Fetched from redis")
			var user pb.UserResponse
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				return &user, nil
			}
			// fallthrough to DB if unmarshal fails
		}
	}

	// Fetch from Postgres
	log.Println("Fetching from DB")
	user := &pb.UserResponse{}
	err := r.db.QueryRow(ctx,
		"SELECT id, name, email FROM users WHERE id=$1", id,
	).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Cache for next time
	if r.redis != nil {
		data, _ := json.Marshal(user)
		_ = r.redis.Set(ctx, "user:"+id, data, 5*time.Minute)
	}

	return user, nil
}