package session

import (
	"clomingo/internal/user"
	"time"
)

type Session struct {
	Id           int64 `datastore:"-"`
	SessionToken string
	SocialToken  string
	SessionType  SessionType
	UserId       int64
	User         user.User `datastore:"-"`
	UserAgent    string
	UserIp       string
	DeviceId     string
	PushToken    string
	PushEnabled  bool
	Timestamp    time.Time
	Timeupdate   time.Time
}

type SessionType int

const (
	ANONYMOUS = iota
	GOOGLE
	APPLE
)
