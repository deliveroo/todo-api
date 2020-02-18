// Package protocol translates domain objects to their API response format. This
// requires you to be deliberate about exposing attributes of your domain
// objects and allows domain objects to evolve without implicitly affecting
// their response format.
package protocol

import (
	"time"

	"github.com/deliveroo/todo-api/domain"
)

// P helps transform entities to the response protocol.
type P struct {
	Debug bool
}

type Task struct {
	ID          int64      `json:"id"`
	Completed   *time.Time `json:"completed"`
	Created     time.Time  `json:"created"`
	Description string     `json:"description"`
}

type TasksAffected struct {
	Updated int64 `json:"updated"`
}

type Account struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type AccountLogin struct {
	Token string `json:"token"`
}

func (p P) AccountLogin(token string) AccountLogin {
	return AccountLogin{
		Token: token,
	}
}

func (p P) Account(v *domain.Account) Account {
	return Account{
		ID:       v.ID,
		Username: v.Username,
	}
}

func (p P) Task(v *domain.Task) Task {
	return Task{
		ID:          v.ID,
		Completed:   v.Completed,
		Created:     v.Created,
		Description: v.Description,
	}
}

func (p P) Tasks(vv []*domain.Task) []Task {
	result := make([]Task, 0, len(vv))
	for _, v := range vv {
		result = append(result, p.Task(v))
	}
	return result
}

func (p P) TasksAffected(v int64) TasksAffected {
	return TasksAffected{
		Updated: v,
	}
}
