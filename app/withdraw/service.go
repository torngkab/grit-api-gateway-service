package withdraw

import (
	"context"
)

type WithdrawService interface {
	Withdraw(ctx context.Context, request WithdrawRequest) (string, error)
}

type withdrawService struct {
	adapter WithdrawAdapter
}

func NewWithdrawService(adapter WithdrawAdapter) WithdrawService {
	return &withdrawService{adapter: adapter}
}

func (s *withdrawService) Withdraw(ctx context.Context, request WithdrawRequest) (string, error) {
	return s.adapter.Withdraw(ctx, request)
}
