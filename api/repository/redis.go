package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"gmsess/domain"
	"strings"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RedisRepository struct {
	redisCli *redis.Client
}

func NewRedisRepository(redisCli *redis.Client) domain.RedisRepository {
	return &RedisRepository{redisCli}
}

func (s *RedisRepository) New(ctx context.Context, sess domain.Session) error {
	bytes, err := json.Marshal(sess)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to marshal session struct with err: %s", err.Error()))
	}

	err = s.redisCli.Set(ctx, sess.SID, bytes, 0).Err()
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to save session struct with err: %s", err.Error()))
	}
	return nil
}

func (s *RedisRepository) Fetch(ctx context.Context, token string) (*domain.Session, error) {
	session := new(domain.Session)
	sid := strings.Split(token, ".")[0]

	res, err := s.redisCli.Get(ctx, sid).Result()
	if err == redis.Nil {
		return &domain.Session{}, status.Error(codes.Unauthenticated, "sid not found")
	} else if err != nil {
		return &domain.Session{}, status.Error(codes.Internal, fmt.Sprintf("failed to fetch session struct with err: %s", err.Error()))
	}

	json.Unmarshal([]byte(res), session)
	return session, nil
}

func (s *RedisRepository) Delete(ctx context.Context, sid string) error {
	err := s.redisCli.Del(ctx, sid).Err()
	return err
}
