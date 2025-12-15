package signin

import (
	"context"
)

type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInService interface {
	SignIn(ctx context.Context, request SignInRequest) (string, error)
}

type signInService struct {
	adapter SignInAdapter
}

func NewSignInService(adapter SignInAdapter) SignInService {
	return &signInService{adapter: adapter}
}

func (s *signInService) SignIn(ctx context.Context, request SignInRequest) (string, error) {
	return s.adapter.SignIn(ctx, request)
}
