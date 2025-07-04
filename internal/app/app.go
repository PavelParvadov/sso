package app

import (
	"go.uber.org/zap"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/repository/sqlite"
	"sso/internal/services/auth"
	"time"
)

type App struct {
	GrpcServer *grpcapp.App
}

func New(log *zap.Logger, grpcPort int, storagePath string, TokenTTL time.Duration) *App {

	storage, err := sqlite.NewStorage(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.NewAuthService(log, storage, storage, storage, TokenTTL)

	grpApp := grpcapp.NewApp(grpcPort, log, authService)

	return &App{
		GrpcServer: grpApp,
	}

}
