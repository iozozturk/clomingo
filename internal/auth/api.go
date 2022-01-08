package auth

import (
	"clomingo/internal/api"
	"clomingo/pkg/config"
	"clomingo/pkg/keys"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type GoogleRequest struct {
	Token     string `json:"token"`
	PushToken string `json:"pushToken"`
}

type Handler struct {
	logger  *zap.Logger
	filters *api.Filter
	conf    *config.ApplicationConf
	authSvc *AuthService
}

func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/google", h.filters.CommonUnAuthenticated(h.GoogleLogin()))
	mux.HandleFunc("/auth/signout", h.filters.CommonAuthenticated(h.SignOut()))
}

func (h *Handler) GoogleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var gr GoogleRequest
		err := json.NewDecoder(r.Body).Decode(&gr)
		if err != nil {
			h.logger.Warn("error in google request", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		session, err := h.authSvc.googleLogin(r.Context(), gr.Token, gr.PushToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.logger.Info("google login success", zap.String("token", gr.Token), zap.Any("session", session))
		w.Header().Set("Content-Type", "application/json")
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionId", Value: session.SessionToken, Expires: expiration}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) SignOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("signout user", zap.Any("session", r.Context().Value(keys.Session)))
		h.authSvc.SignOut(r.Context())
		cookie := http.Cookie{Name: "sessionId", Value: "", Expires: time.Now()}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	}
}

func NewHandler(logger *zap.Logger, filter *api.Filter, conf *config.ApplicationConf, service *AuthService) *Handler {
	return &Handler{
		logger:  logger,
		filters: filter,
		conf:    conf,
		authSvc: service,
	}
}
