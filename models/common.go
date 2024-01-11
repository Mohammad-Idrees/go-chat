package models

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	ConnPool *pgxpool.Pool
}

type Redis struct {
	Client *redis.Client
}

type RequestHeaders struct {
	XRequestID string `json:"X-Request-ID"`
}
