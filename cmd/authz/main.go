package main

import (
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/volvo-cars/connect-access-control/docs"
	"github.com/volvo-cars/connect-access-control/internal/app/authz"
	"github.com/volvo-cars/connect-access-control/internal/config"
	"go.elastic.co/ecszap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	// create logger
	logger := NewLogger(cfg)
	// set default logger
	slog.SetDefault(logger)

	// admin.Run(cfg)
	authz.Run(cfg)
}

func NewLogger(cfg *config.Config) *slog.Logger {
	logLevel, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		logLevel = zapcore.InfoLevel
	}

	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, logLevel)

	opt := &zapslog.HandlerOptions{
		AddSource: true,
	}

	logger := slog.New(zapslog.NewHandler(core, opt))

	logger = logger.With(
		slog.String("service.name", cfg.App.Name),
		slog.String("service.version", cfg.App.Version),
		slog.String("service.environment", cfg.App.Environment.String()),
	)

	return logger
}
