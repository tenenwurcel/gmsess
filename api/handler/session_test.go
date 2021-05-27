package handler

import (
	"context"
	"fmt"
	_entity "gmsess/api/entity"
	_repo "gmsess/api/repository"
	"gmsess/config"
	"gmsess/proto"
	"gmsess/utils"
	"testing"
)

var handler *SessionHandler

func TestMain(m *testing.M) {
	utils.SetupCypher()
	utils.SetupVerifier()
	config.SetupRedis()

	sessionRepo := _repo.NewRedisRepository(config.GetRedisCli())
	sessionEntity := _entity.NewSesssionEntity(sessionRepo)
	handler = NewSessionHandler(sessionEntity)
	m.Run()
}

func TestNew(t *testing.T) {
	sid, err := handler.New(context.Background(), &proto.NewRequest{})
	if err != nil {
		t.Error(err)
	}

	fmt.Println(sid.Sid)
}
