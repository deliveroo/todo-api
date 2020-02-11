package selftest

import (
	"testing"

	"github.com/deliveroo/assert-go"
	"github.com/icrowley/fake"
)

func TestCreateAccount(t *testing.T) {
	t.Run("create with no body", func(t *testing.T) {
		resp := (&API{}).Post(t, "/account", m{})
		resp.AssertStatusCode(t, 400)
	})
	t.Run("create", func(t *testing.T) {
		resp := (&API{}).Post(t, "/account", m{
			"username": fake.UserName(),
			"password": fakePassword(),
		})
		resp.AssertStatusCode(t, 200)
	})
}

func TestCreateAndLoginAccount(t *testing.T) {
	username := fake.UserName()
	password := fakePassword()
	t.Run("create", func(t *testing.T) {
		resp := (&API{}).Post(t, "/account", m{
			"username": username,
			"password": password,
		})
		resp.AssertStatusCode(t, 200)
	})
	t.Run("authorized", func(t *testing.T) {
		resp := (&API{}).Post(t, "/account/login", m{
			"username": username,
			"password": password,
		})
		resp.AssertStatusCode(t, 200)
		tokenBytes := []byte(resp.JSONPathString(t, "token"))
		assert.True(t, len(tokenBytes) >= 32)
	})
	t.Run("wrong username", func(t *testing.T) {
		resp := (&API{}).Post(t, "/account/login", m{
			"username": "wrong-" + username,
			"password": password,
		})
		resp.AssertStatusCode(t, 400)
	})
	t.Run("wrong password", func(t *testing.T) {
		resp := (&API{}).Post(t, "/account/login", m{
			"username": username,
			"password": "wrong-" + password,
		})
		resp.AssertStatusCode(t, 400)
	})
}

func TestAccountValidation(t *testing.T) {
	t.Run("bad username", func(t *testing.T) {
		resp := (&API{}).Post(t, "/account", m{
			"username": "",
			"password": fakePassword(),
		})
		resp.AssertStatusCode(t, 400)
		assert.Equal(t, resp.ErrorMessage(t), "username is required")
	})
	t.Run("bad password", func(t *testing.T) {
		resp := (&API{}).Post(t, "/account", m{
			"username": fake.UserName(),
			"password": "",
		})
		resp.AssertStatusCode(t, 400)
		assert.Contains(t, resp.ErrorMessage(t), "password must")
	})
}

func TestGetAccount(t *testing.T) {
	t.Run("unauthed", func(t *testing.T) {
		resp := (&API{}).Get(t, "/account")
		resp.AssertStatusCode(t, 401)
	})
	t.Run("authed", func(t *testing.T) {
		withAccount(t, func(api *API) {
			resp := api.Get(t, "/account")
			resp.AssertStatusCode(t, 200)
			resp.JSONPathEqual(t, "username", api.Username)
		})
	})
}
