package keys

type ContextKey string

const (
	UserAgent ContextKey = "user-agent"
	UserIp    ContextKey = "user-ip"
	DeviceId  ContextKey = "device-id"
	Language  ContextKey = "language"
	Session   ContextKey = "session"
)
