package repo_test

import (
	"context"
	"testing"

	"github.com/deliveroo/assert-go"
	"github.com/deliveroo/todo-api/domain"
	"github.com/deliveroo/todo-api/repo"
)

func TestCreateAndGetAccount(t *testing.T) {
	var (
		db     = getDB(t)
		client = &repo.Client{db.pool}
		ctx    = context.Background()
	)
	defer db.Close()

	none, err := client.GetAccountByUsername(ctx, "")
	assert.Must(t, err)
	assert.Nil(t, none)

	account := &domain.Account{
		Username:       "username",
		PasswordDigest: "password-digest",
		PasswordSalt:   "password-salt",
	}

	created, err := client.CreateAccount(ctx, account)
	assert.Must(t, err)
	assert.True(t, created.ID != 0)
	assert.Equal(t, created.Username, account.Username)
	assert.Equal(t, created.PasswordDigest, account.PasswordDigest)
	assert.Equal(t, created.PasswordSalt, account.PasswordSalt)
	assert.False(t, created.Created.IsZero())

	got, err := client.GetAccountByUsername(ctx, created.Username)
	assert.Must(t, err)
	assert.Equal(t, got, created)
}
