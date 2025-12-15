package signin

import (
	"context"
	"errors"

	"github.com/torngkab/grit-api-gateway-service/config"

	pb "github.com/torngkab/grit-account-service/account"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type SignInAdapter interface {
	SignIn(ctx context.Context, request SignInRequest) (string, error)
}

type signInAdapter struct {
	config config.Config
}

func NewSignInAdapter(config config.Config) SignInAdapter {
	return &signInAdapter{config: config}
}

func (a *signInAdapter) SignIn(ctx context.Context, request SignInRequest) (string, error) {
	conn, err := grpc.Dial(a.config.Service.AccountService, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewAccountClient(conn)
	resp, err := client.SignIn(ctx, &pb.SignInRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		st := status.Convert(err)
		return "", errors.New(st.Message())
	}

	return resp.GetUserId(), nil
}
