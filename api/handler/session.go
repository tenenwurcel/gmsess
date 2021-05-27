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

func (s *SessionHandler) Authenticate(ctx context.Context, sess *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	authResponse := new(proto.AuthenticateResponse)
	err := s.sessionEntity.Authenticate(ctx, sess, authResponse)
	if err != nil {
		return new(proto.AuthenticateResponse), err
	}

	return authResponse, nil
}

func (s *SessionHandler) New(ctx context.Context, _ *proto.NewRequest) (*proto.NewResponse, error) {
	sess := new(proto.NewResponse)
	err := s.sessionEntity.New(ctx, sess)
	if err != nil {
		return sess, err
	}

	return sess, nil
}
