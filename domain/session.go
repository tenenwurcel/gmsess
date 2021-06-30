package domain

import (
	"context"

	proto "github.com/tenenwurcel/gmprotos/session"

	"github.com/google/uuid"
)

type Session struct {
	Permission   string
	SID          string
	DiscordToken string
	State        string
	RefreshToken string
}

func NewSession() Session {
	return Session{
		Permission:   "",
		SID:          uuid.New().String(),
		DiscordToken: "",
		State:        uuid.New().String(),
	}
}

type SessionEntity interface {
	New(ctx context.Context, sess *proto.NewResponse) error
	Authenticate(ctx context.Context, sess *proto.AuthenticateRequest, sid *proto.AuthenticateResponse) error
	Refresh(ctx context.Context, req *proto.RefreshRequest, res *proto.RefreshResponse) error
	Verify(ctx context.Context, req *proto.VerifyRequest, res *proto.VerifyResponse) error
}
