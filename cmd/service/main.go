package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sinfirst/Test-Task-For-Effective-Mobile/config"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/app"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/middleware/logging"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/router"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/storage/database"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	logger := logging.NewLogger()
	conf, err := config.NewConfig()
	if err != nil {
		logger.Fatalw("can't init config", err)
	}

	db := database.NewPGDB(conf, logger)
	a := app.NewApp(db, conf)
	router := router.NewRouter(a)
	if conf.DatabaseDsn != "" {
		err := database.InitMigrations(conf, logger)
		if err != nil {
			logger.Fatalw("can't init migrations", err)
		}
	}

	server := &http.Server{Addr: conf.ServerAddress, Handler: router}
	go func() {
		logger.Infow("Starting http server", "addr", conf.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalw("create server error: ", err)
		}
	}()

	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Errorw("Server shutdown error", err)
	}
}
