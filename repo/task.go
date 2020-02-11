package repo

import (
	"context"
	"time"

	"github.com/deliveroo/todo-api/domain"
	"github.com/jackc/pgx/v4"
)

// CreateTask inserts a task into the database.
func (c *Client) CreateTask(ctx context.Context, t *domain.Task) (*domain.Task, error) {
	row := c.queryRow(ctx, `
		INSERT INTO tasks (account_id, description, completed)
		VALUES ($1, $2, $3)
		RETURNING id, account_id, description, created, completed;
	`, t.AccountID, t.Description, t.Completed)
	var result domain.Task
	if err := row.Scan(
		&result.ID,
		&result.AccountID,
		&result.Description,
		&result.Created,
		&result.Completed,
	); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteTask deletes a task from the database.
func (c *Client) DeleteTaskByIDAndAccountID(ctx context.Context, taskID, accountID int64) error {
	tag, err := c.exec(ctx, `
		DELETE FROM tasks
		WHERE id = $1
		AND account_id = $2;
	`, taskID, accountID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// UpdateTask updates a task in the database.
func (c *Client) UpdateTask(ctx context.Context, t *domain.Task) (*domain.Task, error) {
	row := c.queryRow(ctx, `
		UPDATE tasks
		SET description = $2, completed = $3
		WHERE id = $1
		RETURNING id, account_id, description, created, completed;
	`, t.ID, t.Description, t.Completed)
	var result domain.Task
	if err := row.Scan(
		&result.ID,
		&result.AccountID,
		&result.Description,
		&result.Created,
		&result.Completed,
	); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTaskByIDAndAccountID fetches a task by ID and account from the database,
// or returns nil if not found.
func (c *Client) GetTaskByIDAndAccountID(ctx context.Context, taskID, accountID int64) (*domain.Task, error) {
	row := c.queryRow(ctx, `
		SELECT id, account_id, description, created, completed
		FROM tasks
		WHERE id = $1
		AND account_id = $2;
	`, taskID, accountID)
	var result domain.Task
	if err := row.Scan(
		&result.ID,
		&result.AccountID,
		&result.Description,
		&result.Created,
		&result.Completed,
	); err != nil {
		if isErrNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// MarkIncompleteTasksCompleteByAccountID marks all incomplete tasks for an account complete.
func (c *Client) MarkIncompleteTasksCompleteByAccountID(ctx context.Context, accountID int64) (int64, error) {
	tag, err := c.exec(ctx, `
		UPDATE tasks
		SET completed = $2
		WHERE account_id = $1
		AND completed IS NULL;
	`, accountID, time.Now().UTC())
	return tag.RowsAffected(), err
}

// GetAllTasksByAccountID fetches all tasks by account from the database.
func (c *Client) GetAllTasksByAccountID(ctx context.Context, accountID int64) ([]*domain.Task, error) {
	rows, err := c.query(ctx, `
		SELECT id, account_id, description, created, completed
		FROM tasks
		WHERE account_id = $1
		ORDER BY created DESC;
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*domain.Task
	for rows.Next() {
		var t domain.Task
		_ = rows.Scan(
			&t.ID,
			&t.AccountID,
			&t.Description,
			&t.Created,
			&t.Completed,
		)
		result = append(result, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
