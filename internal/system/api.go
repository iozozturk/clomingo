package system

import (
	"clomingo/internal/api"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	logger  *zap.Logger
	filters *api.Filter
}

func (h *Handler) HealthCheck() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("health check")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.filters.Logger(h.HealthCheck()))
}

func NewHandler(logger *zap.Logger, filter *api.Filter) *Handler {
	return &Handler{
		logger:  logger,
		filters: filter,
	}
}
