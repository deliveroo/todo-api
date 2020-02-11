package repo

import (
	"context"

	"github.com/deliveroo/todo-api/domain"
)

// CreateAccount inserts an account into the database.
func (c *Client) CreateAccount(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	row := c.queryRow(ctx, `
		INSERT INTO accounts (username, password_digest, password_salt)
		VALUES ($1, $2, $3)
		RETURNING id, username, password_digest, password_salt, created;
	`, a.Username, a.PasswordDigest, a.PasswordSalt)
	var result domain.Account
	if err := row.Scan(
		&result.ID,
		&result.Username,
		&result.PasswordDigest,
		&result.PasswordSalt,
		&result.Created,
	); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAccountByUsername fetches an account by username from the database or
// returns nil if not found.
func (c *Client) GetAccountByUsername(ctx context.Context, username string) (*domain.Account, error) {
	row := c.queryRow(ctx, `
		SELECT id, username, password_digest, password_salt, created
		FROM accounts
		WHERE username = $1;
	`, username)
	var result domain.Account
	if err := row.Scan(
		&result.ID,
		&result.Username,
		&result.PasswordDigest,
		&result.PasswordSalt,
		&result.Created,
	); err != nil {
		if isErrNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// GetAccountByID fetches an account by id from the database or returns nil if
// not found.
func (c *Client) GetAccountByID(ctx context.Context, id int64) (*domain.Account, error) {
	row := c.queryRow(ctx, `
		SELECT id, username, password_digest, password_salt, created
		FROM accounts
		WHERE id = $1;
	`, id)
	var result domain.Account
	if err := row.Scan(
		&result.ID,
		&result.Username,
		&result.PasswordDigest,
		&result.PasswordSalt,
		&result.Created,
	); err != nil {
		if isErrNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
