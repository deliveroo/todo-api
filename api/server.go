// Package api is the API server backed by jsonrest-go. It contains routing,
// handlers, and middleware.
package api

import (
	"net/http"

	"github.com/deliveroo/jsonrest-go"
	"github.com/deliveroo/todo-api/api/protocol"
	"github.com/deliveroo/todo-api/repo"
	"github.com/deliveroo/todo-api/service/session"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Config is the server configuration and dependencies.
type Config struct {
	Database *pgxpool.Pool
	Sessions *session.Service

	DumpErrors bool // render full error in response
}

// Server is an API server.
type Server struct {
	cfg      *Config
	protocol protocol.P
	router   *jsonrest.Router
}

// NewServer configures a new API server.
func NewServer(cfg *Config) *Server {
	s := &Server{
		cfg:      cfg,
		protocol: protocol.P{},
	}
	s.router = router(s)
	return s
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

// Protocol returns the response protocol helper.
func (s *Server) Protocol() protocol.P {
	return s.protocol
}

// Repo returns the repo client.
func (s *Server) Repo() *repo.Client {
	return repo.NewClient(s.cfg.Database)
}

// Sessions returns the sessions service.
func (s *Server) Sessions() *session.Service {
	return s.cfg.Sessions
}
