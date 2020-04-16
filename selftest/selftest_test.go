package selftest

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/deliveroo/todo-api/cmd/todo-api/apicmd"
	"github.com/deliveroo/todo-api/conf"
	"github.com/deliveroo/todo-api/selftest/deps/postgres"
	"github.com/deliveroo/todo-api/selftest/deps/redis"
)

var url string

func TestMain(m *testing.M) {
	if flag.Parse(); testing.Short() {
		return // skip in short mode
	}

	// Connect to Postgres.
	must(postgres.Connect(), "could not connect to postgres")

	// Connect to Redis.
	must(redis.Connect(), "could not connect to redis")

	// Reset the databases.
	must(postgres.Reset(), "error resetting database")
	must(redis.Reset(), "error resettiig redis")

	var addr string
	addr, url = tempAddr()

	cfg := &conf.Config{
		Addr:                addr,
		DatabaseConnTimeout: 5 * time.Second,
		DatabaseMaxConn:     10,
		DatabaseURL:         postgres.URL(),
		Debug:               true,
		MaxSessionDuration:  1 * time.Minute,
		RedisURL:            redis.URL(),
		SuppressLogging:     true,
	}

	// Configure and start API server.
	api, err := apicmd.New(cfg)
	must(err, "error calling apicmd.New")
	go mustDo(api.Run, "error closing server")

	// Wait for API server to listen.
	must(waitForURL(url+"/ping"), "error waiting for API server")

	// Run tests.
	result := m.Run()

	// Shutdown services.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	must(api.Shutdown(ctx), "error shutting down api server")

	os.Exit(result)
}

// must calls log.Fatal if the error is non-nil.
func must(err error, msg string) {
	if err != nil {
		log.Fatalln(msg + ": " + err.Error())
	}
}

// mustDo calls fn, and writes to log.Fatal if an error is returned.
func mustDo(fn func() error, msg string) {
	must(fn(), msg)
}

// tempAddr generates a temporary address on which a TCP server can listen.
func tempAddr() (addr, url string) {
	s := httptest.NewUnstartedServer(nil)
	addr = s.Listener.Addr().String()
	s.Close()
	return addr, "http://" + addr
}

// waitForURL waits for the given URL to respond.
func waitForURL(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		_, err := http.Get(url)
		if err == nil {
			return nil
		}
		select {
		case <-time.After(100 * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
