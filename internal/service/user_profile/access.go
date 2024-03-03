package user_profile_service

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log/slog"
)

type AccessService struct {
	log *slog.Logger
}

func NewAccessService(log *slog.Logger) *AccessService {
	return &AccessService{log: log}
}

func (s *AccessService) GetApiKey(ctx context.Context) string {
	const fn = "user_profile.service.getApiKey"
	log := s.log.With(slog.String("fn", fn))

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		log.Error("metadata is empty")
		return ""
	}
	apiKey, ok := md["apikey"]

	if !ok {
		log.Error("apiKey is empty")
		return ""
	}

	return apiKey[0]
}
