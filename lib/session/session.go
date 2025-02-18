package session

import (
	"context"
	"strconv"
	"time"

	"github.com/bhbdev/jam/lib/logger"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	rdb *redis.Client
}

type UserSession struct {
	Id int `redis:"id"`
}

func NewSessionHandler(redisClient *redis.Client) *Session {
	return &Session{
		rdb: redisClient,
	}
}

func (s *Session) Create(ctx context.Context, sessionId string, userData UserSession, expire time.Duration) error {
	return s.rdb.HSet(ctx, sessionId, userData).Err()
}

func (s *Session) Get(ctx context.Context, sessionId string) (UserSession, error) {
	rs, err := s.rdb.HGetAll(ctx, sessionId).Result()
	if err != nil {
		return UserSession{}, err
	}
	logger.Info("session data", "data", rs)
	userId, err := strconv.Atoi(rs["id"])
	if err != nil {
		return UserSession{}, err
	}
	return UserSession{
		Id: userId,
	}, nil
}

func (s *Session) Delete(ctx context.Context, sessionId string) error {
	return s.rdb.Del(ctx, sessionId).Err()
}

func (s *Session) Expire(ctx context.Context, sessionId string, expire time.Duration) error {
	return s.rdb.Expire(ctx, sessionId, expire).Err()
}
