package repo_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/deliveroo/assert-go"
	"github.com/deliveroo/todo-api/domain"
	"github.com/deliveroo/todo-api/repo"
)

func TestCreateTask(t *testing.T) {
	var (
		db     = getDB(t)
		client = &repo.Client{db.pool}
		ctx    = context.Background()
		now    = time.Now().UTC()
	)
	defer db.Close()
	task := domain.Task{
		AccountID:   1,
		Description: "alpha",
		Completed:   &now,
	}
	result, err := client.CreateTask(ctx, &task)
	assert.Must(t, err)
	assert.Equal(t, task.Description, result.Description)
	assert.Equal(t, task.Completed.Truncate(time.Second), result.Completed.Truncate(time.Second))
	assert.True(t, result.ID != 0)
	assert.False(t, result.Created.IsZero())
}

func TestUpdateTask(t *testing.T) {
	var (
		db     = getDB(t)
		client = &repo.Client{db.pool}
		ctx    = context.Background()
		now    = time.Now().UTC()
	)
	defer db.Close()
	task := &domain.Task{
		AccountID:   1,
		Description: "alpha",
		Completed:   &now,
	}
	task, err := client.CreateTask(ctx, task)
	assert.Must(t, err)
	task.Description = "bravo"
	task.Completed = nil
	{
		updated, err := client.UpdateTask(ctx, task)
		assert.Must(t, err)
		assert.Equal(t, task.Description, updated.Description)
		assert.Nil(t, updated.Completed)
	}
}

func TestGetTaskByIDAndAccountID(t *testing.T) {
	var (
		db     = getDB(t)
		client = &repo.Client{db.pool}
		ctx    = context.Background()
	)

	none, err := client.GetTaskByIDAndAccountID(ctx, 0, 0)
	assert.Must(t, err)
	assert.Nil(t, none)

	task := &domain.Task{
		AccountID:   2,
		Description: "bravo",
		Completed:   nil,
	}
	result, err := client.CreateTask(ctx, task)
	assert.Must(t, err)
	got, err := client.GetTaskByIDAndAccountID(ctx, result.ID, result.AccountID)
	assert.Must(t, err)
	assert.Equal(t, result, got)
	defer db.Close()
}

func TestGetAllTasksByAccountID(t *testing.T) {
	var (
		db        = getDB(t)
		client    = &repo.Client{db.pool}
		ctx       = context.Background()
		accountID = int64(3)
	)
	defer db.Close()
	for i := 0; i < 10; i++ {
		task := domain.Task{
			AccountID:   accountID,
			Description: strconv.Itoa(i),
		}
		_, err := client.CreateTask(ctx, &task)
		assert.Must(t, err)
	}

	tasks, err := client.GetAllTasksByAccountID(ctx, accountID)
	assert.Must(t, err)
	assert.Equal(t, len(tasks), 10)
}

func TestMarkIncompleteTasksCompleteByAccountID(t *testing.T) {
	var (
		db        = getDB(t)
		client    = &repo.Client{db.pool}
		ctx       = context.Background()
		accountID = int64(4)
	)
	defer db.Close()
	for i := 0; i < 10; i++ {
		task := domain.Task{
			AccountID:   accountID,
			Description: strconv.Itoa(i),
			Completed:   nil,
		}
		_, err := client.CreateTask(ctx, &task)
		assert.Must(t, err)
	}
	count, err := client.MarkIncompleteTasksCompleteByAccountID(ctx, accountID)
	assert.Must(t, err)
	assert.Equal(t, count, int64(10))
	tasks, err := client.GetAllTasksByAccountID(ctx, accountID)
	assert.Must(t, err)
	for _, tt := range tasks {
		assert.NotNil(t, tt.Completed)
	}
}

func TestDeleteTask(t *testing.T) {
	var (
		db     = getDB(t)
		client = &repo.Client{db.pool}
		ctx    = context.Background()
		now    = time.Now().UTC()
	)
	defer db.Close()
	task := domain.Task{
		AccountID:   1,
		Description: "alpha",
		Completed:   &now,
	}
	result, err := client.CreateTask(ctx, &task)
	assert.Must(t, err)
	assert.Must(t, client.DeleteTaskByIDAndAccountID(ctx, result.ID, result.AccountID))
	got, err := client.GetTaskByIDAndAccountID(ctx, result.ID, result.AccountID)
	assert.Must(t, err)
	assert.Nil(t, got)
}
