package api

import (
	"context"
	"fmt"

	"github.com/deliveroo/jsonrest-go"
)

func router(s *Server) *jsonrest.Router {
	r := jsonrest.NewRouter()
	r.DumpErrors = s.cfg.DumpErrors

	r.Use(PanicRecoveryMiddleware())

	// Unauthenticated routes.
	unauthed := r.Group()
	unauthed.Routes(jsonrest.RouteMap{
		"POST /account":       s.createAccount,
		"POST /account/login": s.login,
	})

	// Authenticated routes.
	authed := r.Group()
	authed.Use(AuthMiddleware(s))
	authed.Routes(jsonrest.RouteMap{
		// Accounts
		"GET  /account": s.getAccount,

		// Tasks
		"GET    /tasks":     s.getAllTasks,
		"DELETE /tasks/:id": s.deleteTask,
		"GET    /tasks/:id": s.getTask,
		"PUT    /tasks/:id": s.updateTask,
		"POST   /tasks":     s.createTask,
	})

	return r
}

type requestAccountKey struct{}

// AuthMiddleware handles account authentication. If a request isn't
// authenticated, the endpoint handler is not called.
func AuthMiddleware(s *Server) jsonrest.Middleware {
	return func(next jsonrest.Endpoint) jsonrest.Endpoint {
		return func(ctx context.Context, req *jsonrest.Request) (interface{}, error) {
			token := req.Header("x-todo-token")
			sess, err := s.Sessions().Get(ctx, token)
			if err != nil {
				return nil, jsonrest.Unauthorized("unauthorized")
			}
			account, err := s.Repo().GetAccountByID(ctx, sess.AccountID)
			if err != nil {
				return nil, err
			}
			req.Set(requestAccountKey{}, account)
			return next(ctx, req)
		}
	}
}

// PanicRecoveryMiddleware catches and returns any panics that occur in the
// endpoint.
func PanicRecoveryMiddleware() jsonrest.Middleware {
	return func(next jsonrest.Endpoint) jsonrest.Endpoint {
		return func(ctx context.Context, req *jsonrest.Request) (result interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					if e, ok := r.(error); ok {
						err = fmt.Errorf("panic: %w", e)
					} else {
						err = fmt.Errorf("panic: %v", r)
					}
				}
			}()
			result, err = next(ctx, req)
			return
		}
	}
}
