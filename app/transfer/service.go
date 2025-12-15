package transfer

import (
	"context"
)

type TransferService interface {
	Transfer(ctx context.Context, request TransferRequest) (string, error)
}

type transferService struct {
	adapter TransferAdapter
}

func NewTransferService(adapter TransferAdapter) TransferService {
	return &transferService{adapter: adapter}
}

func (s *transferService) Transfer(ctx context.Context, request TransferRequest) (string, error) {
	return s.adapter.Transfer(ctx, request)
}
