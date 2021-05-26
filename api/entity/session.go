package entity

import (
	"errors"
	"fmt"
	"gmsess/domain"
	"gmsess/proto"
	"gmsess/utils"

	"golang.org/x/net/context"
)

type SID struct {
	Sid string
}

type SessionEntity struct {
	repo   domain.SessionRepository
	cypher utils.Cypher
}

func NewSesssionEntity(repo domain.SessionRepository) domain.SessionEntity {
	return &SessionEntity{repo: repo, cypher: utils.GetCypher()}
}

func (e *SessionEntity) New(ctx context.Context, sess *proto.Session) (err error) {
	//Create session
	session := domain.NewSession()

	//Save to storage engine and return encrypted token
	err = e.repo.New(ctx, session)
	if err != nil {
		return
	}

	sess.State = session.State
	sess.Sid, err = e.cypher.Encrypt(session.SID)
	return
}

func (e *SessionEntity) Authenticate(ctx context.Context, sess *proto.Session) (string, error) {
	// Decrypt token
	decryptedSID, err := e.cypher.Decrypt(sess.Sid)
	if err != nil {
		return "", errors.New("Invalid token.")
	}

	// Check if sid exists in the storage engine
	//@TODO Change it to Fetch and check if null in this context
	verified, err := e.repo.Check(ctx, decryptedSID)
	if err != nil {
		return "", errors.New("Internal error.")
	}
	if !verified {
		return "", errors.New("Invalid token.")
	}

	//Get current validator and append it to token - Return sid
	currentValidationToken := utils.Verifier.GetCurrent()
	validatedToken := fmt.Sprintf("%s.%s", decryptedSID, currentValidationToken)
	s, err := e.cypher.Encrypt(validatedToken)
	if err != nil {
		return "", errors.New("Internal error.")
	}

	return s, nil
}

//@TODO Add CheckState function(ctx context.Context, sess *proto.Session) (err error) {}
