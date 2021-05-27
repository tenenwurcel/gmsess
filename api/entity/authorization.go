package entity

import (
	"context"
	"errors"
	"gmsess/domain"
	"gmsess/proto"
	"gmsess/utils"

	"github.com/go-redis/redis/v8"
)

type AuthorizationEntity struct {
	repo   domain.RedisRepository
	cypher utils.Cypher
}

func NewAuthorizationEntity(repo domain.RedisRepository) domain.AuthorizationEntity {
	return &AuthorizationEntity{repo: repo, cypher: utils.GetCypher()}
}

func (e *AuthorizationEntity) Verify(ctx context.Context, verifyReq *proto.VerifyRequest, verifyRes *proto.VerifyResponse) error {
	// Decrypt token
	decryptedSID, err := e.cypher.Decrypt(verifyReq.Sid)
	if err != nil {
		return errors.New("Invalid token.")
	}

	sess, err := e.repo.Fetch(ctx, decryptedSID)
	if err == redis.Nil {
		return errors.New("Invalid token.")
	} else if err != nil {
		return err
	}

	sess.DiscordToken = "" // Delete this line
	//@TODO - Add permission (int) field to session struct
	if verifyReq.WantedPermission != 0 /*sess.Permission*/ {
		return errors.New("Permission denied.")
	}

	return nil
}

func (e *AuthorizationEntity) Refresh(ctx context.Context, refreshReq *proto.RefreshRequest, refreshRes *proto.RefreshResponse) error {
	return nil
}
