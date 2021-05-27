package domain

import (
	"context"
	"gmsess/proto"
)

type AuthorizationEntity interface {
	Verify(ctx context.Context, req *proto.VerifyRequest, res *proto.VerifyResponse) error
	Refresh(ctx context.Context, req *proto.RefreshRequest, res *proto.RefreshResponse) error
}
