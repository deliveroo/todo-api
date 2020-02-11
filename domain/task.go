package domain

import "time"

// Task is a single todo item which belongs to an account.
type Task struct {
	// ID is the database id for the task.
	ID int64

	// AccountID is the database foreign key to the account.
	AccountID int64

	// Completed is the time when the task was marked completed.
	Completed *time.Time

	// Created is the time when the task was created.
	Created time.Time

	// Description is the task description.
	Description string
}
