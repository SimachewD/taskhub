package worker

import (
	"context"
	"log"
	"time"

	"github.com/SimachewD/taskhub/internal/cache"
)

type Worker struct {
	redis *cache.RedisClient
}

func NewWorker(redis *cache.RedisClient) *Worker {
	return &Worker{redis: redis}
}

func (w *Worker) Start(ctx context.Context) {

	for {

		job, err := w.redis.Dequeue(ctx, "task_deleted")

		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		w.handleTaskDeleted(job)
	}
}

func (w *Worker) handleTaskDeleted(taskID string) {

	log.Println("Processing deleted task:", taskID)

	// here we could:
	// send email
	// update analytics
	// clean related data
}