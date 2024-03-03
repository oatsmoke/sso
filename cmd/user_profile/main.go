package main

import (
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"sso/internal/app/user_profile"
	"sso/internal/lib/logger"
	"sso/internal/lib/postgres"
	"syscall"
)

func main() {
	log := logger.SetupLogger()

	if err := godotenv.Load(".env"); err != nil {
		log.Error(err.Error())
		return
	}

	DB, err := postgres.SetupConnection(os.Getenv("DB_STRING"))

	if err != nil {
		log.Error(err.Error())
		return
	}

	server := user_profile_app.NewUserProfileApp(log, DB, os.Getenv("USER_PROFILE_PORT"), os.Getenv("API_KEY"))

	go func() {
		if err := server.Run(); err != nil {
			log.Error(err.Error())
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	server.Stop()

	if err := postgres.CloseConnection(DB); err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("gRPC server stopped")
}
