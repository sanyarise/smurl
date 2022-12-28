package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sanyarise/smurl/config"
	"github.com/sanyarise/smurl/internal/delivery"
	"github.com/sanyarise/smurl/internal/infrastructure/logger"
	"github.com/sanyarise/smurl/internal/infrastructure/server"
	"github.com/sanyarise/smurl/internal/repository"
	"github.com/sanyarise/smurl/internal/usecase"
	"go.uber.org/zap"
)

func main() {
	log.Printf("Start load configuration\n")

	// Config init
	cfg := config.NewConfig()

	// Logger init
	l := logger.NewLogger(cfg.LogLevel)
	logger := l.Logger

	logger.Info("Configuration successfully load")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Database init
	repository, err := repository.NewSmurlRepository(cfg.DNS, logger)
	if err != nil {
		log.Fatal(err)
	}
	// Interface layer init
	usecase := usecase.NewSmurlUsecase(repository, logger)

	// Router init
	router := delivery.NewRouter(usecase, logger, cfg.ServerURL)

	// Server init
	server := server.NewServer(":"+cfg.Port, router, logger, cfg.ReadTimeout, cfg.WriteTimeout, cfg.WriteHeaderTimeout)

	// Start server
	server.Start()
	logger.Info("Start server successfull",
		zap.String("Port ", ":"+cfg.Port))

	<-ctx.Done()

	// Stopping the server when receiving a context termination signal
	server.Stop()
	logger.Info("Server stopped successfull")
	cancel()

	// Database shutdown
	repository.Close()
}
