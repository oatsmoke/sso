package user_profile_app

import (
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log/slog"
	"net"
	"sso/internal/handler/user_profile"
	"sso/internal/service/user_profile"
)

type UserProfileApp struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func NewUserProfileApp(log *slog.Logger, db *gorm.DB, port, apiKey string) *UserProfileApp {
	gRPCServer := grpc.NewServer()

	service := user_profile_service.NewUserProfileService(log, db)
	access := user_profile_service.NewAccessService(log)

	user_profile_handler.NewUserProfileHandler(gRPCServer, service, access, apiKey)

	return &UserProfileApp{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *UserProfileApp) Run() error {
	const fn = "user_profile.app.Run"
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

func (a *UserProfileApp) Stop() {
	const fn = "user_profile.app.Stop"
	log := a.log.With(slog.String("fn", fn))

	log.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
