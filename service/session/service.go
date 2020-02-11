// Package session manages sessions in Redis.
package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Service is the session service, which manages sessions in Redis.
type Service struct {
	Redis              *redis.Pool
	MaxSessionDuration time.Duration
}

// Session stores session data.
type Session struct {
	AccountID int64
}

// New creates and persists a new session.
func (s *Service) New(ctx context.Context, sess *Session) (string, error) {
	key := newSessionID()
	conn, err := s.Redis.GetContext(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	b, err := json.Marshal(sess)
	if err != nil {
		return "", err
	}
	ms := s.MaxSessionDuration.Milliseconds()
	if _, err := conn.Do("SET", key, string(b), "PX", ms); err != nil {
		return "", err
	}
	return key, nil
}

// Get fetches an existing session by token, if it exists or hasn't expired.
func (s *Service) Get(ctx context.Context, token string) (*Session, error) {
	if len(token) < 32 {
		return nil, errors.New("invalid token")
	}
	conn, err := s.Redis.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	reply, err := conn.Do("GET", token)
	if err != nil {
		return nil, err
	}
	v, _ := redis.String(reply, nil)
	var sess Session
	if err := json.Unmarshal([]byte(v), &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

func newSessionID() string {
	return base64.RawURLEncoding.EncodeToString(randomBytes(64))
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}
