package repository

import (
	"context"
	"encoding/json"
	"gmsess/domain"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	redisCli *redis.Client
}

func NewRedisRepository(redisCli *redis.Client) domain.RedisRepository {
	return &RedisRepository{redisCli}
}

func (s *RedisRepository) Fetch(ctx context.Context, sid string) (sess *domain.Session, err error) {
	res, err := s.redisCli.Get(ctx, sid).Result()
	if err == redis.Nil {
		return &domain.Session{}, nil
	} else if err != nil {
		return &domain.Session{}, err
	}

	json.Unmarshal([]byte(res), sess)
	return
}

func (s *RedisRepository) New(ctx context.Context, sess domain.Session) (err error) {
	bytes, err := json.Marshal(sess)
	if err != nil {
		return err
	}

	err = s.redisCli.Set(ctx, sess.SID, bytes, 0).Err()
	return
}
