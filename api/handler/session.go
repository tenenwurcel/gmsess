package handler

import (
	"gmsess/domain"

	proto "github.com/tenenwurcel/gmprotos/session"

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

func (s *SessionHandler) Refresh(ctx context.Context, refreshReq *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	refreshRes := new(proto.RefreshResponse)

	err := s.sessionEntity.Refresh(ctx, refreshReq, refreshRes)
	if err != nil {
		return &proto.RefreshResponse{}, err
	}

	return refreshRes, nil
}

func (s *SessionHandler) Verify(ctx context.Context, verifyReq *proto.VerifyRequest) (*proto.VerifyResponse, error) {
	verifyRes := new(proto.VerifyResponse)

	err := s.sessionEntity.Verify(ctx, verifyReq, verifyRes)
	if err != nil {
		return &proto.VerifyResponse{}, err
	}

	return verifyRes, nil
}
