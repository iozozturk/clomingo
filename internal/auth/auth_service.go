package auth

import (
	"clomingo/internal/session"
	"clomingo/internal/user"
	"clomingo/pkg/config"
	"clomingo/pkg/keys"
	"context"
	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
	"time"
)

type AuthService struct {
	logger         *zap.Logger
	conf           *config.ApplicationConf
	userRepo       *user.UserRepo
	sessionService *session.SessionService
}

func (s *AuthService) googleLogin(ctx context.Context, token string, pushToken string) (*session.Session, error) {
	payload, err := idtoken.Validate(ctx, token, s.conf.GoogleClientIds[2])
	if err != nil {
		s.logger.Warn("cannot google login, wrong id token", zap.Error(err), zap.String("token", token))
		return nil, err
	}
	claims := payload.Claims
	now := time.Now()
	u := &user.User{
		Id:         0,
		Name:       claims["name"].(string),
		Email:      claims["email"].(string),
		Photo:      claims["picture"].(string),
		Timestamp:  now,
		Timeupdate: now,
	}
	existingUser, err := s.userRepo.GetUserByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		userSaved, err := s.userRepo.Save(ctx, u)
		if err != nil {
			return nil, err
		}
		newSession, err := s.sessionService.SaveNewSession(ctx, userSaved.Id, token, pushToken, session.GOOGLE)
		if err != nil {
			return nil, err
		}
		s.logger.Info("google login success", zap.Any("claim", claims), zap.Any("user", existingUser))
		newSession.User = *userSaved
		return newSession, nil
	} else {
		newSession, err := s.sessionService.SaveNewSession(ctx, existingUser.Id, token, pushToken, session.GOOGLE)
		if err != nil {
			return nil, err
		}
		s.logger.Info("google login success", zap.Any("claim", claims), zap.Any("user", existingUser))
		newSession.User = *existingUser
		return newSession, nil
	}
}

func (s *AuthService) SignOut(ctx context.Context) {
	sessionId := ctx.Value(keys.Session).(session.Session).Id
	s.sessionService.Delete(ctx, sessionId)
}

func NewAuthService(logger *zap.Logger, conf *config.ApplicationConf, userRepo *user.UserRepo, sessionService *session.SessionService) *AuthService {
	return &AuthService{
		logger:         logger,
		conf:           conf,
		userRepo:       userRepo,
		sessionService: sessionService,
	}
}
