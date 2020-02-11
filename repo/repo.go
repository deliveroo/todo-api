package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Client provides access to all supported database interactions.
type Client struct {
	Database *pgxpool.Pool
}

// NewClient returns a new repo client.
func NewClient(pool *pgxpool.Pool) *Client {
	return &Client{Database: pool}
}

// queryRow executes the provided query as a prepared statement.
func (c *Client) queryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.Database.QueryRow(ctx, sql, args...)
}

// query executes the provided query as a prepared statement.
func (c *Client) query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.Database.Query(ctx, sql, args...)
}

// exec executes the provided query.
func (c *Client) exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return c.Database.Exec(ctx, sql, args...)
}

func isErrNoRows(err error) bool {
	for err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}
