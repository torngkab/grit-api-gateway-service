package withdraw

import (
	"context"
	"errors"

	"github.com/torngkab/grit-api-gateway-service/config"

	pb "github.com/torngkab/subledger-service/subledger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type WithdrawAdapter interface {
	Withdraw(ctx context.Context, request WithdrawRequest) (string, error)
}

type withdrawAdapter struct {
	config config.Config
}

func NewWithdrawAdapter(config config.Config) WithdrawAdapter {
	return &withdrawAdapter{config: config}
}

func (a *withdrawAdapter) Withdraw(ctx context.Context, request WithdrawRequest) (string, error) {
	conn, err := grpc.Dial(a.config.Service.SubledgerService, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewSubledgerClient(conn)
	resp, err := client.Withdraw(ctx, &pb.WithdrawRequest{
		AccountId: request.AccountId,
		Amount:    float32(request.Amount),
		Note:      request.Note,
	})
	if err != nil {
		st := status.Convert(err)
		return "", errors.New(st.Message())
	}

	return resp.GetTransactionId(), nil
}
