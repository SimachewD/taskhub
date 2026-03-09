package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, dsn string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse postgres config: %v", err)
	}

	// optional tuning
	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.HealthCheckPeriod = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to create postgres pool: %v", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Println("Postgres connected successfully")
	return pool
}