package handler

import (
	"context"
	"gmsess/domain"
	"gmsess/proto"
)

type AuthorizationHandler struct {
	proto.UnimplementedAuthorizationServer
}

func NewAuthorizationHandler(entity domain.AuthorizationEntity) *AuthorizationHandler {
	return &AuthorizationHandler{}
}

func (h *AuthorizationHandler) Verify(ctx context.Context, verifyRequest *proto.VerifyRequest) (*proto.VerifyResponse, error) {
	verifyResponse := new(proto.VerifyResponse)

	return verifyResponse, nil
}

func (h *AuthorizationHandler) Refresh(ctx context.Context, refreshRequest *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	refreshResponse := new(proto.RefreshResponse)

	return refreshResponse, nil
}
