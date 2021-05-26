package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"gmsess/domain"

	"github.com/go-redis/redis/v8"
)

type SessionRepository struct {
	redisCli *redis.Client
}

func NewSessionRepository(redisCli *redis.Client) domain.SessionRepository {
	return &SessionRepository{redisCli}
}

func (s *SessionRepository) Check(ctx context.Context, sid string) (verified bool, err error) {
	//Set default value
	verified = true

	//Check if exists in storage. If not, return false.
	_, err = s.redisCli.Get(ctx, sid).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return
}

func (s *SessionRepository) New(ctx context.Context, sess domain.Session) (err error) {
	bytes, err := json.Marshal(sess)
	fmt.Println(string(bytes))
	fmt.Println("aqui1")
	if err != nil {
		return err
	}
	//Insert session into storage engine and sets it's sid (string) as key
	fmt.Println("aqui2")
	err = s.redisCli.Set(ctx, sess.SID, bytes, 0).Err()
	fmt.Println("aqui3")
	if err != nil {
		return err
	}
	fmt.Println("aqui4")
	return
}
