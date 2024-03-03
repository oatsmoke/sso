package auth_app

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"sso/internal/client/auth"
	"sso/internal/handler/auth"
	"sso/internal/service/auth"
	"time"
)

type AuthApp struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func NewAuthApp(log *slog.Logger, port string, client *auth_client.UserProfileClient, tokenTTL time.Duration, secret string) *AuthApp {
	gRPCServer := grpc.NewServer()

	authService := auth_service.NewAuthService(log, tokenTTL, secret)

	auth_handler.NewAuthHandler(gRPCServer, authService, client)

	return &AuthApp{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *AuthApp) Run() error {
	const fn = "auth.app.Run"
	log := a.log.With(slog.String("fn", fn))

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))

	if err != nil {
		return fmt.Errorf("%s:%w", fn, err)
	}

	log.Info("gRPC server is listening", slog.String("port", a.port))

	if err := a.gRPCServer.Serve(listen); err != nil {
		return fmt.Errorf("%s:%w", fn, err)
	}

	return nil
}

func (a *AuthApp) Stop() {
	const fn = "auth.app.Stop"
	log := a.log.With(slog.String("fn", fn))

	log.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
