package createaccount

import (
	"context"
)

type CreateUserService interface {
	CreateUser(ctx context.Context, request CreateUserRequest) (string, error)
}

type createUserService struct {
	adapter CreateUserAdapter
}

func NewCreateUserService(adapter CreateUserAdapter) CreateUserService {
	return &createUserService{adapter: adapter}
}

func (s *createUserService) CreateUser(ctx context.Context, request CreateUserRequest) (string, error) {
	return s.adapter.CreateUser(ctx, request)
}
