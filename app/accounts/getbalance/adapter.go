package getbalance

import (
	"context"

	"github.com/torngkab/grit-api-gateway-service/config"

	pb "github.com/torngkab/grit-account-service/account"

	"google.golang.org/grpc"
)

type GetBalanceAdapter interface {
	GetBalance(ctx context.Context, request GetBalanceRequest) (GetBalanceResponse, error)
}

type getBalanceAdapter struct {
	config config.Config
}

func NewGetBalanceAdapter(config config.Config) GetBalanceAdapter {
	return &getBalanceAdapter{config: config}
}

func (a *getBalanceAdapter) GetBalance(ctx context.Context, request GetBalanceRequest) (GetBalanceResponse, error) {
	conn, err := grpc.Dial(a.config.Service.AccountService, grpc.WithInsecure())
	if err != nil {
		return GetBalanceResponse{}, err
	}
	defer conn.Close()

	client := pb.NewAccountClient(conn)
	resp, err := client.GetAccountBalance(ctx, &pb.GetAccountBalanceRequest{
		AccountId: request.AccountId,
	})
	if err != nil {
		return GetBalanceResponse{}, err
	}

	return GetBalanceResponse{
		Balance:         float64(resp.GetBalance()),
		LatestUpdatedAt: resp.GetLatestUpdatedAt(),
	}, nil
}
