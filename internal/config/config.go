package config

import (
	"fmt"

	"github.com/SimachewD/taskhub/internal/utils"
)

type Config struct {
	PostgresURL string
	RedisAddr   string
}

func Load() *Config {
	pgHost := utils.GetEnv("POSTGRES_HOST", "localhost")
	pgPort := utils.GetEnv("POSTGRES_PORT", "55432")
	pgUser := utils.GetEnv("POSTGRES_USER", "postgres")
	pgPass := utils.GetEnv("POSTGRES_PASSWORD", "postgres")
	pgDB := utils.GetEnv("POSTGRES_DB", "taskhub")

	redisAddr := utils.GetEnv("REDIS_ADDR", "localhost:6379")

	return &Config{
		PostgresURL: fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			pgUser, pgPass, pgHost, pgPort, pgDB,
		),
		RedisAddr: redisAddr,
	}
}