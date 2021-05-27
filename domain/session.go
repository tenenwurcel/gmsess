package domain

import (
	"context"
	"gmsess/proto"

	"github.com/google/uuid"
)

type Session struct {
	SID          string
	DiscordToken string
	State        string
	RefreshToken string
}

func NewSession() Session {
	return Session{
		SID:          uuid.New().String(),
		DiscordToken: "",
		State:        uuid.New().String(),
	}
}

type SessionEntity interface {
	New(ctx context.Context, sess *proto.NewResponse) error
	Authenticate(ctx context.Context, sess *proto.AuthenticateRequest, sid *proto.AuthenticateResponse) error
}
