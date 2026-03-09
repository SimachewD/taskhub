package config

import (
	"fmt"
	"os"
)

type Config struct {
	PostgresURL string
	RedisAddr   string
}

func Load() *Config {
	pgHost := getEnv("POSTGRES_HOST", "localhost")
	pgPort := getEnv("POSTGRES_PORT", "55432")
	pgUser := getEnv("POSTGRES_USER", "postgres")
	pgPass := getEnv("POSTGRES_PASSWORD", "postgres")
	pgDB := getEnv("POSTGRES_DB", "taskhub")

	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")

	return &Config{
		PostgresURL: fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			pgUser, pgPass, pgHost, pgPort, pgDB,
		),
		RedisAddr: redisAddr,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}