package createaccount

import (
	"context"

	"github.com/torngkab/grit-api-gateway-service/config"

	pb "github.com/torngkab/grit-account-service/account"

	"google.golang.org/grpc"
)

type CreateUserAdapter interface {
	CreateUser(ctx context.Context, request CreateUserRequest) (string, error)
}

type createUserAdapter struct {
	config config.Config
}

func NewCreateUserAdapter(config config.Config) CreateUserAdapter {
	return &createUserAdapter{config: config}
}

func (a *createUserAdapter) CreateUser(ctx context.Context, request CreateUserRequest) (string, error) {
	conn, err := grpc.Dial(a.config.Service.AccountService, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewAccountClient(conn)
	resp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Username:     request.Username,
		Password:     request.Password,
		ReferralCode: request.ReferralCode,
	})
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}
