package api

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/deliveroo/jsonrest-go"
	"github.com/deliveroo/todo-api/domain"
	"github.com/deliveroo/todo-api/service/session"
)

type accountParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p accountParams) validate() error {
	if len(p.Username) == 0 {
		return errors.New("username is required")
	}
	if utf8.RuneCountInString(p.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

// login is POST /account/login
func (s *Server) login(ctx context.Context, req *jsonrest.Request) (interface{}, error) {
	var params accountParams
	if err := req.BindBody(&params); err != nil {
		return nil, err
	}
	if err := params.validate(); err != nil {
		return nil, jsonrest.BadRequest(err.Error())
	}
	account, err := s.Repo().GetAccountByUsername(ctx, params.Username)
	if err != nil {
		return nil, err
	}
	ok, err := account.Authenticate(params.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, jsonrest.BadRequest("incorrect username or password")
	}
	sess := session.Session{
		AccountID: account.ID,
	}
	token, err := s.Sessions().New(ctx, &sess)
	if err != nil {
		return nil, err
	}
	return s.Protocol().AccountLogin(token), nil
}

// getAccount is GET /account
func (s *Server) getAccount(ctx context.Context, req *jsonrest.Request) (interface{}, error) {
	account := req.Get(requestAccountKey{}).(*domain.Account)
	return s.Protocol().Account(account), nil
}

// createAccount is POST /account
func (s *Server) createAccount(ctx context.Context, req *jsonrest.Request) (interface{}, error) {
	var params accountParams
	if err := req.BindBody(&params); err != nil {
		return nil, err
	}
	if err := params.validate(); err != nil {
		return nil, jsonrest.BadRequest(err.Error())
	}
	account := &domain.Account{
		Username: params.Username,
	}
	if err := account.SetPassword(params.Password); err != nil {
		return nil, err
	}
	account, err := s.Repo().CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	return s.Protocol().Account(account), nil
}
