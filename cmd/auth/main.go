package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"sso/internal/app/auth"
	"sso/internal/client/auth"
	"sso/internal/lib/logger"
	"syscall"
	"time"
)

func main() {
	log := logger.SetupLogger()

	if err := godotenv.Load(".env"); err != nil {
		log.Error(err.Error())
		return
	}

	adds := fmt.Sprintf("%s:%s", os.Getenv("USER_PROFILE_ADDR"), os.Getenv("USER_PROFILE_PORT"))

	client, err := auth_client.NewUserProfileClient(context.Background(), adds, os.Getenv("API_KEY"))

	if err != nil {
		log.Error(err.Error())
		return
	}

	tokenTTL, err := time.ParseDuration(os.Getenv("TOKEN_TTL"))

	if err != nil {
		log.Error(err.Error())
		return
	}

	server := auth_app.NewAuthApp(log, os.Getenv("AUTH_PORT"), client, tokenTTL, os.Getenv("SECRET"))

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

	log.Info("gRPC server stopped")
}
