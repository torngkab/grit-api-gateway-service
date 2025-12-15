package transfer

import (
	"context"
	"errors"

	"github.com/torngkab/grit-api-gateway-service/config"

	pb "github.com/torngkab/subledger-service/subledger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type TransferAdapter interface {
	Transfer(ctx context.Context, request TransferRequest) (string, error)
}

type transferAdapter struct {
	config config.Config
}

func NewTransferAdapter(config config.Config) TransferAdapter {
	return &transferAdapter{config: config}
}

func (a *transferAdapter) Transfer(ctx context.Context, request TransferRequest) (string, error) {
	conn, err := grpc.Dial(a.config.Service.SubledgerService, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewSubledgerClient(conn)
	resp, err := client.Transfer(ctx, &pb.TransferRequest{
		FromAccountId: request.FromAccountId,
		ToAccountId:   request.ToAccountId,
		Amount:        float32(request.Amount),
		Note:          request.Note,
	})
	if err != nil {
		st := status.Convert(err)
		return "", errors.New(st.Message())
	}

	return resp.GetTransactionId(), nil
}
