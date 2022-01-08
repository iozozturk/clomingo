package api

import (
	"clomingo/internal/session"
	"clomingo/pkg/keys"
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Filter struct {
	logger         *zap.Logger
	sessionService *session.SessionService
}

func (f *Filter) CommonAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return f.Logger(f.Validated(f.ExtractMeta(f.Authenticated(next))))
}
func (f *Filter) CommonUnAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return f.Logger(f.Validated(f.ExtractMeta(next)))
}

func (f *Filter) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		maybeSession, _ := r.Cookie("sessionId")
		defer f.logger.Info("request processed",
			zap.Duration("duration", time.Now().Sub(startTime)),
			zap.String("device-id", r.Header.Get("X-Installation-Id")),
			zap.String("user-ip", r.RemoteAddr),
			zap.String("session-id", maybeSession.String()),
			zap.String("user-agent", r.UserAgent()))
		next(w, r)
	}
}

func (f *Filter) Authenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("sessionId")
		if err != nil {
			f.logger.Info("request does not have sessionId", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Request does not carry session info"))
			return
		}
		sessionFound := f.sessionService.GetSessionAndUserByToken(r.Context(), sessionCookie.Value)
		if sessionFound == nil {
			f.logger.Info("session with sessionId not found", zap.String("sessionId", sessionCookie.Value))
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("session expired"))
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), keys.Session, *sessionFound))
		next(w, r)
	}
}

func (f *Filter) Validated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		installId := r.Header.Get("X-Installation-Id")
		if installId == "" {
			f.logger.Info("request is not from valid source", zap.String("error", "no installation id"))
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Request is not from valid source"))
		} else {
			next(w, r)
		}
	}
}

func (f *Filter) ExtractMeta(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), keys.UserAgent, r.UserAgent()))
		r = r.WithContext(context.WithValue(r.Context(), keys.UserIp, r.RemoteAddr))
		r = r.WithContext(context.WithValue(r.Context(), keys.DeviceId, r.Header.Get("X-Installation-Id")))
		r = r.WithContext(context.WithValue(r.Context(), keys.Language, r.Header.Get("Accept-Language")))
		next(w, r)
	}
}

func NewFilter(logger *zap.Logger, service *session.SessionService) *Filter {
	return &Filter{
		logger:         logger,
		sessionService: service,
	}
}
