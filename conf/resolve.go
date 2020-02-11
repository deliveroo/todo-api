package conf

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/deliveroo/todo-api/service/session"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Dependencies are the resolved dependencies.
type Dependencies struct {
	Database  *pgxpool.Pool
	RedisPool *redis.Pool
	Sessions  *session.Service
}

// Resolve resolves the application dependencies using its config.
func Resolve(ctx context.Context, c *Config) (*Dependencies, error) {
	db, err := resolveDatabase(ctx, c)
	if err != nil {
		return nil, err
	}

	redisPool, err := resolveRedisPool(c)
	if err != nil {
		return nil, fmt.Errorf("redisPool: %w", err)
	}

	return &Dependencies{
		Database:  db,
		RedisPool: redisPool,
		Sessions: &session.Service{
			Redis:              redisPool,
			MaxSessionDuration: c.MaxSessionDuration,
		},
	}, nil
}

func resolveDatabase(ctx context.Context, c *Config) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(c.DatabaseURL)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = c.DatabaseMaxConn
	cfg.MaxConnLifetime = c.DatabaseConnTimeout
	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	_, err = pool.Exec(ctx, "select 1;")
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func resolveRedisPool(c *Config) (*redis.Pool, error) {
	if c.RedisURL == "" {
		return nil, errors.New("RedisURL is required")
	}
	pool := &redis.Pool{
		MaxActive: c.RedisMaxActive,
		MaxIdle:   c.RedisMaxIdle,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(c.RedisURL,
				// Read timeout on server should be greater than ping period.
				redis.DialReadTimeout(45*time.Second),
				redis.DialWriteTimeout(10*time.Second),
			)
		},
	}
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	return pool, nil
}
