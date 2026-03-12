package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SimachewD/taskhub/internal/cache"
	pb "github.com/SimachewD/taskhub/proto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db    *pgxpool.Pool
	redis *cache.RedisClient
}

func NewTaskRepository(db *pgxpool.Pool, redis *cache.RedisClient) *TaskRepository {
	return &TaskRepository{db: db, redis: redis}
}

func (r *TaskRepository) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {

	id := uuid.NewString()

	_, err := r.db.Exec(ctx,
		`INSERT INTO tasks (id,title,description,user_id)
		 VALUES ($1,$2,$3,$4)`,
		id,
		req.Title,
		req.Description,
		req.UserId,
	)

	if err != nil {
		return nil, err
	}

	task := &pb.TaskResponse{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		UserId:      req.UserId,
		Completed:   false,
	}

	// cache task
	data, _ := json.Marshal(task)
	_ = r.redis.Set(ctx, "task:"+id, data, 5*time.Minute)

	return task, nil
}

func (r *TaskRepository) GetTask(ctx context.Context, id string) (*pb.TaskResponse, error) {

	// Redis first
	cached, err := r.redis.Get(ctx, "task:"+id)

	if err == nil {
		var task pb.TaskResponse
		json.Unmarshal([]byte(cached), &task)
		return &task, nil
	}

	task := &pb.TaskResponse{}
	
	var createdAt time.Time // temporary variable for DB timestamp

	err = r.db.QueryRow(ctx,
		`SELECT id,title,description,user_id,completed,created_at
		 FROM tasks WHERE id=$1`, id).
		Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.UserId,
			&task.Completed,
			&createdAt, // scan timestamp here
		)

	if err != nil {
		fmt.Printf("Error fetching task: %v", err)
		return nil, err
	}

	// convert time.Time to string
    task.CreatedAt = createdAt.Format(time.RFC3339)

	data, _ := json.Marshal(task)
	_ = r.redis.Set(ctx, "task:"+id, data, 5*time.Minute)

	return task, nil
}

func (r *TaskRepository) ListTasks(
	ctx context.Context,
	userID string,
	limit int32,
	cursor string,
) ([]*pb.TaskResponse, string, error) {

	if cursor == "" {
		cursor = "infinity" //default pointer
	}
	if limit <= 0 {
		limit = 20 // default page size
	}

	query := `
	SELECT id, title, description, user_id, completed, created_at
	FROM tasks
	WHERE user_id = $1
		AND created_at < $2::timestamptz
	ORDER BY created_at DESC
	LIMIT $3
	`
	args := []any{userID, cursor, limit}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}

	defer rows.Close()

	var tasks []*pb.TaskResponse
	var lastCursor string

	for rows.Next() {

		task := &pb.TaskResponse{}

		var createdAt time.Time // temporary variable for DB timestamp
		
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.UserId,
			&task.Completed,
			&createdAt, // scan timestamp here
		)

		if err != nil {
			return nil, "", err
		}

		// convert time.Time to string
    	task.CreatedAt = createdAt.Format(time.RFC3339)
		lastCursor = task.CreatedAt

		tasks = append(tasks, task)
	}

	return tasks, lastCursor, nil
}