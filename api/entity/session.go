package entity

import (
	"fmt"
	"gmsess/config"
	"gmsess/domain"
	"gmsess/proto"
	"gmsess/utils"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SID struct {
	Sid string
}

type SessionEntity struct {
	repo   domain.RedisRepository
	cypher utils.Cypher
}

func NewSesssionEntity(repo domain.RedisRepository) domain.SessionEntity {
	return &SessionEntity{repo: repo, cypher: config.GetCypher()}
}

func (e *SessionEntity) New(ctx context.Context, sess *proto.NewResponse) (err error) {
	//Create session
	session := domain.NewSession()

	//Save to storage engine and return encrypted token
	err = e.repo.New(ctx, session)
	if err != nil {
		return
	}

	sess.Expire = config.Verifier.GetExpiration()
	sess.RefreshToken, err = e.cypher.Encrypt(session.RefreshToken)
	sess.State = session.State
	sess.Sid, err = e.cypher.Encrypt(session.SID)
	return
}

func (e *SessionEntity) Authenticate(ctx context.Context, sess *proto.AuthenticateRequest, sid *proto.AuthenticateResponse) error {
	// Decrypt token
	decryptedSID, err := e.cypher.Decrypt(sess.Sid)
	if err != nil {
		return err
	}

	// Check if sid exists
	curSess, err := e.repo.Fetch(ctx, decryptedSID)

	if err != nil {
		return err
	} else if curSess.State != sess.State {
		e.repo.Delete(ctx, sess.Sid)
		return status.Error(codes.Unauthenticated, "CSRF validation failed")
	}

	//Validate permission
	if sess.Permission < "a" || sess.Permission > "c" {
		return status.Error(codes.InvalidArgument, "permission must be between 0 and 2")
	}

	//Sets permission - Redis overwrite the existing pair
	curSess.Permission = sess.Permission
	err = e.repo.New(ctx, *curSess)
	if err != nil {
		return err
	}

	//Get current validator and append it to token - Return sid
	currentValidationToken := config.Verifier.GetCurrent()
	authenticatedToken := fmt.Sprintf("%s.%s", decryptedSID, currentValidationToken)
	sid.Sid, err = e.cypher.Encrypt(authenticatedToken)
	if err != nil {
		return err
	}

	return nil
}

func (e *SessionEntity) Verify(ctx context.Context, verifyReq *proto.VerifyRequest, verifyRes *proto.VerifyResponse) error {
	// Decrypt token
	decryptedSID, err := e.cypher.Decrypt(verifyReq.Sid)
	if err != nil {
		return err
	}

	//Split decrypted token - sid.token
	splitedToken := strings.Split(decryptedSID, ".")
	if len(splitedToken) < 2 {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	//Check if token verifier matches current token
	err = config.Verifier.Verify(splitedToken[1])
	if err != nil {
		return err
	}

	//Find current session
	sess, err := e.repo.Fetch(ctx, decryptedSID)
	if err != nil {
		return err
	}

	//Check if user has permission
	if verifyReq.WantedPermission > sess.Permission {
		return status.Error(codes.Unauthenticated, "permission denied")
	}

	verifyRes.Valid = true
	return nil
}

func (e *SessionEntity) Refresh(ctx context.Context, refreshReq *proto.RefreshRequest, refreshRes *proto.RefreshResponse) error {
	// Decrypt token
	decryptedSID, err := e.cypher.Decrypt(refreshReq.Sid)
	if err != nil {
		return err
	}

	//Split decrypted token - sid.token
	splitedToken := strings.Split(decryptedSID, ".")
	if len(splitedToken) < 2 {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	//Check if token verifier matches current token
	err = config.Verifier.Verify(splitedToken[1])
	if err != nil {
		grpcError := status.Convert(err)

		if grpcError.Code() == codes.PermissionDenied {
			decryptedToken := fmt.Sprintf("%s.%s", splitedToken[0], config.Verifier.GetCurrent())

			token, err := e.cypher.Encrypt(decryptedToken)

			if err != nil {
				return err
			}

			refreshRes.Sid = token
			return nil
		}

		return err
	}

	refreshRes.Sid = refreshReq.Sid
	return nil
}
