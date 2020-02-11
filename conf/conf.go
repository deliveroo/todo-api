// Package conf holds all of the configuration for the application: e.g.
// database connection strings, port to listen on, external credentials,
// environment. By locating all configuration in one place, it's easy to see all
// parameters at a glance. This approach implies that no other packages should
// access environment variables directly.
package conf

import (
	"time"

	"github.com/caarlos0/env/v6"
)

// Config is the configuration needed to bootstrap the application's
// dependencies.
type Config struct {
	Addr                string        `env:"ADDR" envDefault:":4000"`                       // Server listen address
	DatabaseConnTimeout time.Duration `env:"DATABASE_CONN_TIMEOUT" envDefault:"10s"`        // Postgres connection timeout
	DatabaseMaxConn     int32         `env:"DATABASE_MAX_CONN" envDefault:"10"`             // Postgres connection pool limit
	DatabaseURL         string        `env:"DATABASE_URL"`                                  // Postgres connection string
	Debug               bool          `env:"DEBUG"`                                         // Enable debug mode
	MaxSessionDuration  time.Duration `env:"MAX_SESSION_DURATION" envDefault:"24h"`         // The maximum duration of a login session.
	RedisMaxActive      int           `env:"REDIS_MAX_ACTIVE" envDefault:"5"`               // Max active redis pool connections
	RedisMaxIdle        int           `env:"REDIS_MAX_IDLE" envDefault:"5"`                 // Maximum idle redis pool connections
	RedisURL            string        `env:"REDIS_URL" envDefault:"redis://127.0.0.1:6379"` // Redis connection string
	SuppressLogging     bool          `env:"SUPPRESS_LOGGING"`                              // Suppress logging, useful for testing
}

// Load loads the application configuration from command line flags and
// environment variables.
func Load() (*Config, error) {
	c := Config{}
	if err := env.Parse(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
