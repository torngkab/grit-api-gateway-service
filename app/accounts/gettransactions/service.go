package gettransactions

import (
	"context"

	pbAccount "github.com/torngkab/grit-account-service/account"
	pb "github.com/torngkab/subledger-service/subledger"
)

type GetTransactionsService interface {
	GetTransactions(ctx context.Context, request GetTransactionsRequest) (pb.GetTransactionsResponse, error)
	GetAccountIdsByUserId(ctx context.Context, userId string) ([]*pbAccount.AccountModel, error)
}

type getTransactionsService struct {
	adapter GetTransactionsAdapter
}

func NewGetTransactionsService(adapter GetTransactionsAdapter) GetTransactionsService {
	return &getTransactionsService{adapter: adapter}
}

func (s *getTransactionsService) GetTransactions(ctx context.Context, request GetTransactionsRequest) (pb.GetTransactionsResponse, error) {
	return s.adapter.GetTransactions(ctx, request)
}

func (s *getTransactionsService) GetAccountIdsByUserId(ctx context.Context, userId string) ([]*pbAccount.AccountModel, error) {
	return s.adapter.GetAccountIdsByUserId(ctx, userId)
}
