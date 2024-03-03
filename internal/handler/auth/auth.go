package auth_handler

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sso/internal/client/auth"
	ssov1 "sso/internal/gen/sso/v1"
)

type AuthService interface {
	GetToken(ctx context.Context, userId int64) (string, error)
}

type AuthHandler struct {
	ssov1.UnimplementedAuthServer
	authService AuthService
	client      *auth_client.UserProfileClient
}

func NewAuthHandler(gRPC *grpc.Server, authService AuthService, client *auth_client.UserProfileClient) {
	ssov1.RegisterAuthServer(gRPC, &AuthHandler{
		authService: authService,
		client:      client,
	})
}

func (h *AuthHandler) Registration(ctx context.Context, req *ssov1.RegistrationRequest) (*ssov1.RegistrationResponse, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	userId, err := h.client.Create(ctx, req.GetLogin(), req.GetPassword(), req.GetEmail())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	if userId == 0 {
		return nil, status.Error(codes.AlreadyExists, "already exists")
	}

	return &ssov1.RegistrationResponse{UserId: userId}, nil
}

func (h *AuthHandler) SignIn(ctx context.Context, req *ssov1.SignInRequest) (*ssov1.SignInResponse, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userId, err := h.client.Authentication(ctx, req.GetLogin(), req.GetPassword())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	if userId == 0 {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	token, err := h.authService.GetToken(ctx, userId)

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.SignInResponse{Token: token}, nil
}

func (h *AuthHandler) GetUser(ctx context.Context, req *ssov1.GetUserRequest) (*ssov1.GetUserResponse, error) {

	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	user, err := h.client.GetUser(ctx, req.GetUserId())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	if user == nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &ssov1.GetUserResponse{User: user}, nil
}
