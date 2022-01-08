package main

import (
	"clomingo/internal"
	"clomingo/pkg/config"
	"clomingo/pkg/monitoring"
	"clomingo/pkg/server"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	logger := monitoring.NewLogger().With(zap.String("app", "clomingo"))
	defer monitoring.Sync(logger)
	logger.Info("starting clomingo backend")
	appConf := config.NewApplication()
	mux := http.NewServeMux()

	internal.NewApi(logger, appConf).InitApis(mux)

	srv := server.New(mux, appConf.Host)
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatal("cannot start application", zap.Error(err))
	}
}
