package gettransactions

import (
	"context"

	"github.com/torngkab/grit-api-gateway-service/config"

	pbAccount "github.com/torngkab/grit-account-service/account"
	pb "github.com/torngkab/subledger-service/subledger"

	"google.golang.org/grpc"
)

type GetTransactionsAdapter interface {
	GetTransactions(ctx context.Context, request GetTransactionsRequest) (pb.GetTransactionsResponse, error)
	GetAccountIdsByUserId(ctx context.Context, userId string) ([]*pbAccount.AccountModel, error)
}

type getTransactionsAdapter struct {
	config config.Config
}

func NewGetTransactionsAdapter(config config.Config) GetTransactionsAdapter {
	return &getTransactionsAdapter{config: config}
}

func (a *getTransactionsAdapter) GetTransactions(ctx context.Context, request GetTransactionsRequest) (pb.GetTransactionsResponse, error) {
	conn, err := grpc.Dial(a.config.Service.SubledgerService, grpc.WithInsecure())
	if err != nil {
		return pb.GetTransactionsResponse{}, err
	}
	defer conn.Close()

	client := pb.NewSubledgerClient(conn)
	resp, err := client.GetTransactions(ctx, &pb.GetTransactionsRequest{
		AccountIds: request.AccountIds,
		Page:       request.Page,
		Limit:      request.Limit,
	})
	if err != nil {
		return pb.GetTransactionsResponse{}, err
	}

	return pb.GetTransactionsResponse{
		Transactions: resp.Transactions,
		Total:        resp.Total,
		Page:         resp.Page,
		Limit:        resp.Limit,
	}, nil
}

func (a *getTransactionsAdapter) GetAccountIdsByUserId(ctx context.Context, userId string) ([]*pbAccount.AccountModel, error) {
	conn, err := grpc.Dial(a.config.Service.AccountService, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbAccount.NewAccountClient(conn)
	resp, err := client.GetAccountsByUserId(ctx, &pbAccount.GetAccountsByUserIdRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetAccounts(), nil
}
