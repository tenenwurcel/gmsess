package entity

import (
	"errors"
	"fmt"
	"gmsess/domain"
	"gmsess/proto"
	"gmsess/utils"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type SID struct {
	Sid string
}

type SessionEntity struct {
	repo   domain.RedisRepository
	cypher utils.Cypher
}

func NewSesssionEntity(repo domain.RedisRepository) domain.SessionEntity {
	return &SessionEntity{repo: repo, cypher: utils.GetCypher()}
}

func (e *SessionEntity) New(ctx context.Context, sess *proto.NewResponse) (err error) {
	//Create session
	session := domain.NewSession()

	//Save to storage engine and return encrypted token
	err = e.repo.New(ctx, session)
	if err != nil {
		return
	}

	sess.State = session.State
	sess.Sid, err = e.cypher.Encrypt(session.SID)
	if err != nil {
		return
	}
	sess.RefreshToken, err = e.cypher.Encrypt(session.RefreshToken)
	return
}

func (e *SessionEntity) Authenticate(ctx context.Context, sess *proto.AuthenticateRequest, sid *proto.AuthenticateResponse) error {
	// Decrypt token
	decryptedSID, err := e.cypher.Decrypt(sess.Sid)
	if err != nil {
		return errors.New("Invalid token.")
	}

	// Check if sid exists in the storage engine
	curSess, err := e.repo.Fetch(ctx, decryptedSID)
	if err == redis.Nil {
		return errors.New("Invalid token.")
	} else if err != nil {
		return err
	} else if curSess.State != sess.State {
		return errors.New("Invalid CRSF token.")
	}

	//Get current validator and append it to token - Return sid
	currentValidationToken := utils.Verifier.GetCurrent()
	authenticatedToken := fmt.Sprintf("%s.%s", decryptedSID, currentValidationToken)
	sid.Sid, err = e.cypher.Encrypt(authenticatedToken)
	if err != nil {
		return errors.New("Internal error.")
	}

	return nil
}

/*func (e *SessionEntity) Authorize(ctx context.Context, sess *proto.Session) error {
	return nil
}*/

//func (e *SessionEntity) Refresh(ctx context.Context, sess *proto.Session)
