package service

import (
	"encoding/json"
	"log"

	"github.com/SimachewD/taskhub/internal/cache"
	pb "github.com/SimachewD/taskhub/proto"
)

type NotificationService struct {
	redis *cache.RedisClient
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationServer(redis *cache.RedisClient) *NotificationService {
	return &NotificationService{redis: redis}
}

func (s *NotificationService) SubscribeTasks(req *pb.SubscribeTasksRequest, stream pb.NotificationService_SubscribeTasksServer) error {

	pubsub := s.redis.Subscribe("tasks_channel:" + req.UserId)
	ch := pubsub.Channel()

	for msg := range ch {
		var payload struct {
			Event string           `json:"event"`
			Task  *pb.TaskResponse `json:"task"`
		}
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			log.Println("Failed to unmarshal task event:", err)
			continue
		}
		if err := stream.Send(&pb.NotificationResponse{
                Event: payload.Event,
                Task:  payload.Task, // directly send the TaskResponse
            }); err != nil {
			return err
		}
	}

	return nil
}