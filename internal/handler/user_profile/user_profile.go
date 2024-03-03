package user_profile_handler

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	ssov1 "sso/internal/gen/sso/v1"
	"sso/internal/model"
)

type UserProfileService interface {
	Create(ctx context.Context, login, password, email string) (int64, error)
	Authentication(ctx context.Context, login, password string) int64
	Info(ctx context.Context, userId int64) (*model.User, error)
}

type AccessService interface {
	GetApiKey(ctx context.Context) string
}

type UserProfileHandler struct {
	ssov1.UnimplementedUserProfileServer
	userProfileService UserProfileService
	accessService      AccessService
	apiKey             string
}

func NewUserProfileHandler(gRPC *grpc.Server, userProfileService UserProfileService, accessService AccessService, apiKey string) {
	ssov1.RegisterUserProfileServer(gRPC, &UserProfileHandler{
		userProfileService: userProfileService,
		accessService:      accessService,
		apiKey:             apiKey,
	})
}

func (h *UserProfileHandler) Create(ctx context.Context, req *ssov1.CreateRequest) (*ssov1.CreateResponse, error) {
	apiKey := h.accessService.GetApiKey(ctx)

	if apiKey != h.apiKey {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	userId, err := h.userProfileService.Create(ctx, req.GetLogin(), req.GetPassword(), req.GetEmail())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.CreateResponse{UserId: userId}, nil
}

func (h *UserProfileHandler) Authentication(ctx context.Context, req *ssov1.AuthenticationRequest) (*ssov1.AuthenticationResponse, error) {
	apiKey := h.accessService.GetApiKey(ctx)

	if apiKey != h.apiKey {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userId := h.userProfileService.Authentication(ctx, req.GetLogin(), req.GetPassword())

	return &ssov1.AuthenticationResponse{UserId: userId}, nil
}

func (h *UserProfileHandler) Info(ctx context.Context, req *ssov1.InfoRequest) (*ssov1.InfoResponse, error) {
	apiKey := h.accessService.GetApiKey(ctx)

	if apiKey != h.apiKey {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	user, err := h.userProfileService.Info(ctx, req.GetUserId())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	if user == nil {
		return nil, nil
	}

	return &ssov1.InfoResponse{User: mapUserToApi(user)}, nil
}

func mapUserToApi(user *model.User) *ssov1.User {
	return &ssov1.User{
		Id:    user.ID,
		Login: user.Login,
		Email: user.Email,
	}
}
