// The apicmd package hoists configuration and application startup together.
// This allows the application entrypoint to be remain lean and makes it easier
// to configure and run the application from selftest for integration testing.
package apicmd

import (
	"context"
	"net/http"

	"github.com/deliveroo/todo-api/api"
	"github.com/deliveroo/todo-api/conf"
	"go.uber.org/zap"
)

// Command orchestrates running and stopping the API server.
type Command struct {
	cancel context.CancelFunc
	dep    *conf.Dependencies
	server *http.Server
}

// New generates a new API server command for the given config.
func New(cfg *conf.Config) (*Command, error) {
	ctx, cancel := context.WithCancel(context.Background())
	dep, err := conf.Resolve(ctx, cfg)
	if err != nil {
		cancel()
		return nil, err
	}
	api := api.NewServer(&api.Config{
		Database:   dep.Database,
		DumpErrors: cfg.Debug,
		Sessions:   dep.Sessions,
	})
	return &Command{
		cancel: cancel,
		dep:    dep,
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: api,
		},
	}, nil
}

// Run starts the API server.
func (c *Command) Run() error {
	zap.L().Info("apicmd.Run", zap.String("addr", c.server.Addr))
	err := c.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown commences graceful shutdown of the API server.
func (c *Command) Shutdown(ctx context.Context) error {
	c.cancel()
	return c.server.Shutdown(ctx)
}
