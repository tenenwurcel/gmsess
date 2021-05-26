package handler

import (
	"gmsess/domain"
	"gmsess/proto"

	"golang.org/x/net/context"
)

type SessionHandler struct {
	sessionEntity domain.SessionEntity
	proto.UnimplementedAuthenticatorServer
}

func NewSessionHandler(sessionEntity domain.SessionEntity) *SessionHandler {
	return &SessionHandler{sessionEntity: sessionEntity}
}

func (s *SessionHandler) Authenticate(ctx context.Context, sess *proto.Session) (sid *proto.SID, err error) {
	sid.S, err = s.sessionEntity.Authenticate(ctx, sess)

	return
}

func (s *SessionHandler) New(ctx context.Context, _ *proto.Void) (*proto.Session, error) {
	sess := new(proto.Session)
	err := s.sessionEntity.New(ctx, sess)
	if err != nil {
		return sess, err
	}

	return sess, nil
}
