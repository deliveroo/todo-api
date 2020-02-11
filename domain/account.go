package domain

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"time"
)

// Account is a user account.
type Account struct {
	// ID is the database id for the account.
	ID int64

	// Created is when the account was created.
	Created time.Time

	// PasswordDigest is the SHA512 digest of the salt combined with the
	// plain-text password.
	PasswordDigest string

	// PasswordSalt is random bytes for securing the password digest.
	PasswordSalt string

	// Username is the account username.
	Username string
}

// SetPassword creates and sets a new password digest and salt.
func (a *Account) SetPassword(password string) error {
	if a == nil {
		return errors.New("account cannot be nil")
	}
	salt := base64.RawURLEncoding.EncodeToString(randomBytes(64))
	digest := sha512.Sum512([]byte(salt + password))
	a.PasswordDigest = base64.RawURLEncoding.EncodeToString(digest[:])
	a.PasswordSalt = salt
	return nil
}

// Authenticate compares an incoming plain-text password with the expected
// password digest. It is safe to use with a nil account.
func (a *Account) Authenticate(password string) (bool, error) {
	var (
		err  error
		want []byte
		ok   = true
		salt string
	)
	if a != nil {
		var derr error
		want, derr = base64.RawURLEncoding.DecodeString(a.PasswordDigest)
		if err != nil {
			ok = false
			err = derr
		}
		salt = a.PasswordSalt
	}
	digest := sha512.Sum512([]byte(salt + password))
	if len(want) < 64 {
		want = randomBytes(64)
		ok = false
	}
	if subtle.ConstantTimeCompare(digest[:], want) == 0 {
		ok = false
	}
	return ok, err
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}
