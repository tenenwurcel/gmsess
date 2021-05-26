package domain

import (
	"context"
	"encoding/json"
	"gmsess/proto"

	"github.com/google/uuid"
)

type Session struct {
	SID          string
	DiscordToken string
	State        string
}

func NewSession() Session {
	return Session{
		SID:          uuid.New().String(),
		DiscordToken: "",
		State:        "",
	}
}

type SessionEntity interface {
	New(ctx context.Context, sess *proto.Session) error
	Authenticate(ctx context.Context, sess *proto.Session) (string, error)
}

type SessionRepository interface {
	New(ctx context.Context, session Session) error
	Check(ctx context.Context, sid string) (verified bool, err error)
}

func (s Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
