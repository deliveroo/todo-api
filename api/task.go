package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/deliveroo/jsonrest-go"
	"github.com/deliveroo/todo-api/domain"
	"github.com/jackc/pgx"
)

type taskParams struct {
	Description string     `json:"description"`
	Completed   *time.Time `json:"completed"`
}

func (p taskParams) validate() error {
	if len(p.Description) == 0 {
		return errors.New("description is required")
	}
	return nil
}

// createTask is POST /tasks
func (s *Server) createTask(ctx context.Context, req *jsonrest.Request) (interface{}, error) {
	account := req.Get(requestAccountKey{}).(*domain.Account)
	var params taskParams
	if err := req.BindBody(&params); err != nil {
		return nil, err
	}
	if err := params.validate(); err != nil {
		return nil, jsonrest.BadRequest(err.Error())
	}
	t := &domain.Task{
		AccountID:   account.ID,
		Description: params.Description,
		Completed:   params.Completed,
	}
	t, err := s.Repo().CreateTask(ctx, t)
	if err != nil {
		return nil, err
	}
	return s.Protocol().Task(t), nil
}

// updateTask is PUT /tasks/:id
func (s *Server) updateTask(ctx context.Context, req *jsonrest.Request) (interface{}, error) {
	account := req.Get(requestAccountKey{}).(*domain.Account)
	var params taskParams
	if err := req.BindBody(&params); err != nil {
		return nil, err
	}
	if err := params.validate(); err != nil {
		return nil, jsonrest.BadRequest(err.Error())
	}
	tid, _ := strconv.ParseInt(req.Param("id"), 10, 64)
	t, err := s.Repo().GetTaskByIDAndAccountID(ctx, tid, account.ID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, jsonrest.NotFound(fmt.Sprintf("task not found, id=%d", tid))
	}
	t.Description = params.Description
	t.Completed = params.Completed
	t, err = s.Repo().UpdateTask(ctx, t)
	if err != nil {
		return nil, err
	}
	return s.Protocol().Task(t), nil
}

// deleteTask is DELETE /tasks/:id
func (s *Server) deleteTask(ctx context.Context, req *jsonrest.Request) (interface{}, error) {
	account := req.Get(requestAccountKey{}).(*domain.Account)
	tid, _ := strconv.ParseInt(req.Param("id"), 10, 64)
	err := s.Repo().DeleteTaskByIDAndAccountID(ctx, tid, account.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, jsonrest.NotFound(fmt.Sprintf("task not found, id=%d", tid))
		}
		return nil, err
	}
	return nil, nil
}

// getTasks is GET /tasks
func (s *Server) getAllTasks(ctx context.Context, req *jsonrest.Request) (interface{}, error) {
	account := req.Get(requestAccountKey{}).(*domain.Account)
	tasks, err := s.Repo().GetAllTasksByAccountID(ctx, account.ID)
	if err != nil {
		return nil, err
	}
	return s.Protocol().Tasks(tasks), nil
}

// getTask is GET /tasks/:id
func (s *Server) getTask(ctx context.Context, req *jsonrest.Request) (interface{}, error) {
	account := req.Get(requestAccountKey{}).(*domain.Account)
	tid, _ := strconv.ParseInt(req.Param("id"), 10, 64)
	t, err := s.Repo().GetTaskByIDAndAccountID(ctx, tid, account.ID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, jsonrest.NotFound(fmt.Sprintf("task not found, id=%d", tid))
	}
	return s.Protocol().Task(t), nil
}

// markAllTasksComplete is POST /tasks/all/complete
func (s *Server) markAllTasksComplete(ctx context.Context, req *jsonrest.Request) (interface{}, error) {
	account := req.Get(requestAccountKey{}).(*domain.Account)
	marked, err := s.Repo().MarkIncompleteTasksCompleteByAccountID(ctx, account.ID)
	return s.Protocol().TasksAffected(marked), err
}
