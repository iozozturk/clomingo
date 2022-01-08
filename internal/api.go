package internal

import (
	"clomingo/internal/api"
	"clomingo/internal/auth"
	"clomingo/internal/session"
	"clomingo/internal/system"
	"clomingo/internal/user"
	"clomingo/pkg/config"
	"clomingo/pkg/db"
	"go.uber.org/zap"
	"net/http"
)

type Api struct {
	logger *zap.Logger
	conf   *config.ApplicationConf
}

func (a *Api) InitApis(mux *http.ServeMux) {
	// init all dependencies and endpoints
	datastore := db.NewDatastore(a.conf, a.logger)
	sessionRepo := session.NewSessionRepo(datastore, a.logger)
	userRepo := user.NewUserRepo(datastore, a.logger)
	sessionService := session.NewSessionService(a.logger, a.conf, sessionRepo, userRepo)
	authSvc := auth.NewAuthService(a.logger, a.conf, userRepo, sessionService)
	filters := api.NewFilter(a.logger, sessionService)

	system.NewHandler(a.logger, filters).SetupRoutes(mux)
	auth.NewHandler(a.logger, filters, a.conf, authSvc).SetupRoutes(mux)
}

func NewApi(logger *zap.Logger, conf *config.ApplicationConf) *Api {
	return &Api{
		logger: logger,
		conf:   conf,
	}
}
