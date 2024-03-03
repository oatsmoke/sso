package auth_client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	ssov1 "sso/internal/gen/sso/v1"
)

type UserProfileClient struct {
	api    ssov1.UserProfileClient
	apiKey string
}

func NewUserProfileClient(ctx context.Context, addr, apiKey string) (*UserProfileClient, error) {
	const fn = "auth.client.NewUserProfileClient"

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &UserProfileClient{
		api:    ssov1.NewUserProfileClient(conn),
		apiKey: apiKey}, nil
}

func (c *UserProfileClient) Create(ctx context.Context, login, password, email string) (int64, error) {
	const fn = "auth.client.Registration"

	ctx = metadata.AppendToOutgoingContext(ctx, "apikey", c.apiKey)

	resp, err := c.api.Create(ctx, &ssov1.CreateRequest{Login: login, Password: password, Email: email})

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return resp.UserId, nil
}

func (c *UserProfileClient) Authentication(ctx context.Context, login, password string) (int64, error) {
	const fn = "auth.client.Authentication"

	ctx = metadata.AppendToOutgoingContext(ctx, "apikey", c.apiKey)

	resp, err := c.api.Authentication(ctx, &ssov1.AuthenticationRequest{Login: login, Password: password})

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return resp.UserId, nil
}

func (c *UserProfileClient) GetUser(ctx context.Context, userId int64) (*ssov1.User, error) {
	const fn = "auth.client.GetUser"

	ctx = metadata.AppendToOutgoingContext(ctx, "apikey", c.apiKey)

	resp, err := c.api.Info(ctx, &ssov1.InfoRequest{UserId: userId})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return resp.User, nil
}
