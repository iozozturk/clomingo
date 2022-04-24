package main

import (
	"clomingo/internal"
	"clomingo/pkg/config"
	"clomingo/pkg/monitoring"
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	logger := monitoring.NewLogger().With(zap.String("app", "clomingo"))
	defer monitoring.Sync(logger)
	logger.Info("starting clomingo backend")
	appConf := config.NewApplication()

	server := internal.NewApi(logger, appConf).InitApis()

	go func() {
		logger.Info("Server started", zap.String("host", appConf.Host))
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Cannot start application", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Could not shutdown server", zap.Error(err))
	}
	logger.Info("Bye bye")
}
