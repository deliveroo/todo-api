package session_test

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/deliveroo/assert-go"
	"github.com/deliveroo/todo-api/selftest/deps/redis"
	"github.com/deliveroo/todo-api/service/session"
)

func TestMain(m *testing.M) {
	if flag.Parse(); testing.Short() {
		return // skip in short mode
	}

	// Connect to Redis.
	must(redis.Connect(), "could not connect to redis")

	// Run tests.
	result := m.Run()

	// Reset the database.
	must(redis.Reset(), "error resetting redis")

	os.Exit(result)
}

func TestSessionPersistence(t *testing.T) {
	var (
		ctx = context.Background()
		s   = &session.Service{
			Redis:              redis.Pool(),
			MaxSessionDuration: 100 * time.Millisecond,
		}
		sess  = &session.Session{AccountID: time.Now().Unix()}
		token string
	)
	t.Run("new", func(t *testing.T) {
		var err error
		token, err = s.New(ctx, sess)
		assert.Must(t, err)
		assert.True(t, token != "")
	})
	t.Run("get", func(t *testing.T) {
		got, err := s.Get(ctx, token)
		assert.Must(t, err)
		assert.Equal(t, sess, got)
	})
	t.Run("expires", func(t *testing.T) {
		time.Sleep(110 * time.Millisecond)
		got, err := s.Get(ctx, token)
		assert.NotNil(t, err)
		assert.Nil(t, got)
	})
}

// must calls log.Fatal if the error is non-nil.
func must(err error, msg string) {
	if err != nil {
		log.Fatalln(msg + ": " + err.Error())
	}
}
