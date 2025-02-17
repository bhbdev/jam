package config

const (
	SessionCookieName     = "session_id"
	SessionCookieMaxAge   = 60 * 60 * 2 // 2 hours
	SessionRedisKeyFormat = "session:%s"
)
