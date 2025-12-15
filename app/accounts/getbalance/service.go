package getbalance

import (
	"context"
)

type GetBalanceService interface {
	GetBalance(ctx context.Context, request GetBalanceRequest) (GetBalanceResponse, error)
}

type getBalanceService struct {
	adapter GetBalanceAdapter
}

func NewGetBalanceService(adapter GetBalanceAdapter) GetBalanceService {
	return &getBalanceService{adapter: adapter}
}

func (s *getBalanceService) GetBalance(ctx context.Context, request GetBalanceRequest) (GetBalanceResponse, error) {
	return s.adapter.GetBalance(ctx, request)
}
