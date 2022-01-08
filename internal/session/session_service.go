package session

import (
	"clomingo/internal/user"
	"clomingo/pkg/config"
	"clomingo/pkg/keys"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type SessionService struct {
	logger   *zap.Logger
	conf     *config.ApplicationConf
	repo     *SessionRepo
	userRepo *user.UserRepo
}

func (s *SessionService) SaveNewSession(ctx context.Context, userId int64, socialToken string, pushToken string, sessionType SessionType) (*Session, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		s.logger.Error("error creating uuid", zap.Error(err))
		return &Session{}, err
	}
	now := time.Now()
	deviceId := ctx.Value(keys.DeviceId).(string)
	session := &Session{
		SessionToken: newUUID.String(),
		SocialToken:  socialToken,
		SessionType:  sessionType,
		UserId:       userId,
		UserAgent:    ctx.Value(keys.UserAgent).(string),
		UserIp:       ctx.Value(keys.UserIp).(string),
		DeviceId:     deviceId,
		PushToken:    pushToken,
		PushEnabled:  pushToken != "",
		Timestamp:    now,
		Timeupdate:   now,
	}
	ids, _ := s.repo.findKeysByDeviceId(ctx, deviceId)
	if ids != nil {
		_ = s.repo.deleteMany(ctx, ids)
	}
	sessionDb, err := s.repo.save(ctx, session)
	if err != nil {
		s.logger.Error("error saving new session", zap.Error(err))
		return nil, err
	}
	return sessionDb, nil
}

func (s *SessionService) GetSessionAndUserByToken(ctx context.Context, sessionToken string) *Session {
	sessionDb, _ := s.repo.GetBySessionToken(ctx, sessionToken)
	if sessionDb == nil {
		return nil
	}
	userDb, _ := s.userRepo.GetByUserId(ctx, sessionDb.UserId)
	if userDb == nil {
		return nil
	}
	sessionDb.User = *userDb
	return sessionDb
}

func (s *SessionService) Delete(ctx context.Context, sessionId int64) {
	s.repo.deleteById(ctx, sessionId)
}

func NewSessionService(logger *zap.Logger, conf *config.ApplicationConf, sessionRepo *SessionRepo, userRepo *user.UserRepo) *SessionService {
	return &SessionService{
		logger:   logger,
		conf:     conf,
		repo:     sessionRepo,
		userRepo: userRepo,
	}
}
