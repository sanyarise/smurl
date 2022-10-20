package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/sanyarise/smurl/config"
	"github.com/sanyarise/smurl/internal/infrastructure/api/handler"
	"github.com/sanyarise/smurl/internal/infrastructure/api/logger"
	"github.com/sanyarise/smurl/internal/infrastructure/api/routeropenapi"
	"github.com/sanyarise/smurl/internal/infrastructure/db/pgstore"
	"github.com/sanyarise/smurl/internal/infrastructure/server"
	"github.com/sanyarise/smurl/internal/usecases/repos/smurlrepo"

	"go.uber.org/zap"
)

func main() {
	log.Printf("start load configuration.\n")

	// Config init
	cfg := config.NewConfig()

	// Logger init
	l := logger.NewLogger(cfg.LogLevel)
	defer l.Logger.Sync()
	logger := l.Logger

	logger.Info("configuration file successfully load.")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	// Database init
	smst, err := pgstore.NewSmurlStore(cfg.DNS, logger)
	if err != nil {
		log.Fatal(err)
	}
	// Interface layer init
	smr := smurlrepo.NewSmurlStorage(smst, logger)

	// Handlers init
	hs := handler.NewHandlers(smr, logger)

	// Router init
	router := routeropenapi.NewRouterOpenAPI(hs, logger, cfg.ServerURL)

	// Server init
	server := server.NewServer(":"+cfg.Port, router, logger, cfg.ReadTimeout, cfg.WriteTimeout, cfg.WriteHeaderTimeout)

	// Start server
	server.Start(smr)
	logger.Info("Start server successfull",
		zap.String("port ", ":"+cfg.Port))

	<-ctx.Done()

	// Stopping the server when receiving a context termination signal
	server.Stop()
	logger.Info("Server stopped successfull")
	cancel()
	
	// Database shutdown
	smst.Close()
}
