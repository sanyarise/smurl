package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/sanyarise/smurl/config"
	"github.com/sanyarise/smurl/internal/infrastructure/api/handler"
	"github.com/sanyarise/smurl/internal/infrastructure/api/routeropenapi"
	"github.com/sanyarise/smurl/internal/infrastructure/db/pgstore"
	"github.com/sanyarise/smurl/internal/infrastructure/server"
	"github.com/sanyarise/smurl/internal/usecases/repos/smurlrepo"

	"go.uber.org/zap"
)

func main() {
	log.Printf("start load configuration.\n")

	//Инициализация конфигурации
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error on load config %v\n", err)
	}

	defer cfg.Logger.Sync()

	l := cfg.Logger

	l.Info("configuration file successfully load.")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	//Инициализация базы данных
	smst, err := pgstore.NewSmurlStore(os.Getenv("DATABASE_URL"), l)
	if err != nil {
		log.Fatal(err)
	}
	//Инициализация слоя с интерфейсами
	smr := smurlrepo.NewSmurlStorage(smst, l)

	//Инициализация хэндлеров
	hs := handler.NewHandlers(smr, l)

	//router := smurlapi.NewRouterChi(hs, l) - инициализация рукописного
	//роутера chi

	//Инициализация сгенерированного роутера chi
	router := routeropenapi.NewRouterOpenAPI(hs, l, cfg.ServerURL)

	//Инициализация сервера
	server := server.NewServer(cfg.Port, router, l, cfg.ReadTimeout, cfg.WriteTimeout, cfg.WriteHeaderTimeout)

	//Запуск сервера
	server.Start(smr)
	l.Info("Start server successfull",
		zap.String("port ", cfg.Port))

	<-ctx.Done()

	//Остановка сервера при получении сигнала о завершении контекста
	server.Stop()
	l.Info("Server stopped successfull")
	cancel()
	//Завершение работы базы данных
	smst.Close()
}
