package domain_test

import (
	"testing"

	"github.com/deliveroo/assert-go"
	"github.com/deliveroo/todo-api/domain"
)

func TestAccountAuthenticate(t *testing.T) {
	t.Run("SetPassword and Authenticate", func(t *testing.T) {
		account := &domain.Account{}
		assert.Must(t, account.SetPassword("very-secret"))
		{
			ok, err := account.Authenticate("very-secret")
			assert.Must(t, err)
			assert.True(t, ok)
		}
		{
			ok, err := account.Authenticate("wrong-password")
			assert.Must(t, err)
			assert.False(t, ok)
		}
	})
	t.Run("when account is nil", func(t *testing.T) {
		var account *domain.Account
		assert.NotNil(t, account.SetPassword("very-secret"))
		ok, err := account.Authenticate("very-secret")
		assert.Must(t, err)
		assert.False(t, ok)
	})
}
