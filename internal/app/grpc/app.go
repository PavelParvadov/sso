package grpcApp

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	AuthGRPC "sso/internal/grpc/auth"
)

type App struct {
	Grpc   *grpc.Server
	Logger *zap.Logger
	Port   int
	AuthGRPC.Auth
}

func NewApp(port int, logger *zap.Logger, authService AuthGRPC.Auth) *App {
	server := grpc.NewServer()
	AuthGRPC.Register(server, authService)
	return &App{
		Grpc:   server,
		Logger: logger,
		Port:   port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Port))
	if err != nil {
		return err
	}
	a.Logger.Info("Starting gRPC server: ", zap.Int("port", a.Port))
	if err := a.Grpc.Serve(l); err != nil {
		return fmt.Errorf("grpc serve: %w", err)
	}
	return nil
}

func (a *App) Stop() {
	a.Logger.Info("Stopping gRPC server: ", zap.Int("port", a.Port))
	a.Grpc.GracefulStop()
	a.Logger.Info("gRPC server stopped gracefully")
}
