// Package redis provides a Redis dependency for selftest.
package redis

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

func URL() string {
	return os.Getenv("TODO_API_TEST_REDIS_URL")
}

func Conn() (redis.Conn, error) {
	return redis.DialURL(URL())
}

func Pool() *redis.Pool {
	pool := &redis.Pool{
		Wait: true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(URL(),
				redis.DialReadTimeout(45*time.Second),
				redis.DialWriteTimeout(10*time.Second),
			)
		},
	}
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		panic(err)
	}
	return pool
}

func Connect() error {
	conn, err := Conn()
	if err != nil {
		return err
	}
	reply, err := conn.Do("PING")
	if err != nil {
		return err
	}
	s, err := redis.String(reply, nil)
	if err != nil {
		return err
	}
	if s != "PONG" {
		return errors.New("unexpected reply")
	}
	return err
}

func Reset() error {
	conn, err := redis.DialURL(URL())
	if err != nil {
		return err
	}
	_, err = conn.Do("FLUSHALL")
	return err
}
