package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/deliveroo/todo-api/conf"
	"github.com/oklog/run"
	"go.uber.org/zap"
)

func main() {
	var (
		logger *zap.Logger
		g      run.Group
	)
	cfg, err := conf.Load()
	if err != nil {
		log.Fatalln(err)
	}

	if cfg.SuppressLogging {
		logger = zap.NewNop()
	} else {
		var err error
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalln(err)
		}
	}
	defer func() {
		_ = logger.Sync()
	}()
	_ = zap.ReplaceGlobals(logger)

	// Signal handler.
	{
		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(error) {
			cancel()
		})
	}

	if err := g.Run(); err != nil {
		zap.L().Fatal("run error", zap.Error(err))
	}
}
