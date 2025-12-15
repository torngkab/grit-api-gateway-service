package deposit

import (
	"context"

	"github.com/torngkab/grit-api-gateway-service/config"

	pb "github.com/torngkab/subledger-service/subledger"

	"google.golang.org/grpc"
)

type DepositAdapter interface {
	Deposit(ctx context.Context, request DepositRequest) (string, error)
}

type depositAdapter struct {
	config config.Config
}

func NewDepositAdapter(config config.Config) DepositAdapter {
	return &depositAdapter{config: config}
}

func (a *depositAdapter) Deposit(ctx context.Context, request DepositRequest) (string, error) {
	conn, err := grpc.Dial(a.config.Service.SubledgerService, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewSubledgerClient(conn)
	resp, err := client.Deposit(ctx, &pb.DepositRequest{
		AccountId: request.AccountId,
		Amount:    float32(request.Amount),
		Note:      request.Note,
	})
	if err != nil {
		return "", err
	}

	return resp.GetTransactionId(), nil
}
