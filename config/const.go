package config

import "time"

const (
	SessionCookieName     = "session_id"
	SessionCookieMaxAge   = 60 * 60 * 2 // 2 hours
	SessionRedisKeyFormat = "session:%s"
	SessionExpireTime     = 8 * time.Hour
)
