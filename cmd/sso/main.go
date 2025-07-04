package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/pkg/logging"
	"syscall"
	"time"
)

const (
	LOCAL = "local"
	DEV   = "dev"
	PROD  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	log := logging.GetLogger()
	log.Info("Старт", zap.Any("env", cfg.Env))
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GrpcServer.MustRun()
	//GS

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GrpcServer.Stop()

	log.Info("application stopped")
	time.Sleep(500 * time.Millisecond)

}

func SetupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case LOCAL:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case DEV:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case PROD:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	return logger
}
