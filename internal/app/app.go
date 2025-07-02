package app

import (
	"go.uber.org/zap"
	grpcapp "sso/internal/app/grpc"
	"time"
)

type App struct {
	GrpcApp *grpcapp.App
}

func New(log *zap.Logger, grpcPort int, storagePath string, TokenTTL time.Duration) *App {

	//TODO: хранилище
	//TODO: сервис auth

	grpcApp := grpcapp.NewApp(grpcPort, log)

	return &App{
		GrpcApp: grpcApp,
	}

}
