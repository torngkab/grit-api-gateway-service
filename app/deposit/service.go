package deposit

import (
	"context"
)

type DepositService interface {
	Deposit(ctx context.Context, request DepositRequest) (string, error)
}

type depositService struct {
	adapter DepositAdapter
}

func NewDepositService(adapter DepositAdapter) DepositService {
	return &depositService{adapter: adapter}
}

func (s *depositService) Deposit(ctx context.Context, request DepositRequest) (string, error) {
	return s.adapter.Deposit(ctx, request)
}
