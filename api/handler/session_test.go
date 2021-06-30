package handler

import (
	"context"
	"encoding/hex"
	_entity "gmsess/api/entity"
	_repo "gmsess/api/repository"
	"gmsess/config"
	"strings"
	"testing"

	proto "github.com/tenenwurcel/gmprotos/session"
)

//@TODO - Create functions to validate tokens - too verbose.
var lHandler *SessionHandler

func TestMain(m *testing.M) {
	config.SetupCypher()
	config.SetupVerifier()
	config.SetupRedis()

	sessionRepo := _repo.NewRedisRepository(config.GetRedisCli())
	sessionEntity := _entity.NewSesssionEntity(sessionRepo)
	lHandler = NewSessionHandler(sessionEntity)
	m.Run()
}

func TestNew(t *testing.T) {
	sid, err := lHandler.New(context.Background(), &proto.NewRequest{})
	if err != nil {
		t.Error(err)
		return
	}
	if len(sid.Sid) != 128 {
		t.Error("invalid sid len")
		return
	}

	cypher := config.GetCypher()

	decryptedToken, err := cypher.Decrypt(sid.Sid)
	if err != nil {
		t.Error(err)
		return
	}

	splitToken := strings.Split(decryptedToken, "-")
	if len(splitToken) != 5 {
		t.Error("invalid UUID")
		return
	}
	for _, tokenPart := range splitToken {
		_, err = hex.DecodeString(tokenPart)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestAuthenticate(t *testing.T) {
	newRes, err := lHandler.New(context.Background(), &proto.NewRequest{})
	if err != nil {
		t.Error(err)
		return
	}

	authReq := &proto.AuthenticateRequest{
		Sid:          newRes.Sid,
		State:        newRes.State,
		Permission:   "c",
		DiscordToken: "1237819287301982730918273091872309817908",
	}
	sid, err := lHandler.Authenticate(context.Background(), authReq)
	if err != nil {
		t.Error(err)
		return
	}
	if len(sid.Sid) != 162 {
		t.Error("invalid authenticated sid len")
	}

	cypher := config.GetCypher()

	decryptedToken, err := cypher.Decrypt(sid.Sid)
	if err != nil {
		t.Error(err)
		return
	}

	splitedToken := strings.Split(decryptedToken, ".")
	if len(splitedToken) != 2 {
		t.Error("invalid authorized token")
		return
	}
	if len(splitedToken[1]) != 16 {
		t.Error("invalid verification token len")
		return
	}

	_, err = hex.DecodeString(splitedToken[1])
	if err != nil {
		t.Error(err)
		return
	}

	decryptedSid := splitedToken[0]
	splitSid := strings.Split(decryptedSid, "-")
	for _, tokenPart := range splitSid {
		_, err = hex.DecodeString(tokenPart)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestVerify(t *testing.T) {
	newRes, err := lHandler.New(context.Background(), &proto.NewRequest{})
	if err != nil {
		t.Error(err)
		return
	}

	authReq := &proto.AuthenticateRequest{
		Sid:          newRes.Sid,
		State:        newRes.State,
		Permission:   "c",
		DiscordToken: "1237819287301982730918273091872309817908",
	}
	sid, err := lHandler.Authenticate(context.Background(), authReq)
	if err != nil {
		t.Error(err)
		return
	}

	verifyReq := &proto.VerifyRequest{
		Sid:              sid.Sid,
		WantedPermission: "c",
	}
	verified, err := lHandler.Verify(context.Background(), verifyReq)
	if err != nil {
		t.Error(err)
		return
	}
	if !verified.Valid {
		t.Error("unable to authorize")
	}
}

func TestRefresh(t *testing.T) {
	cypher := config.GetCypher()
	newRes, err := lHandler.New(context.Background(), &proto.NewRequest{})
	if err != nil {
		t.Error(err)
		return
	}

	authReq := &proto.AuthenticateRequest{
		Sid:          newRes.Sid,
		State:        newRes.State,
		Permission:   "c",
		DiscordToken: "1237819287301982730918273091872309817908",
	}
	sid, err := lHandler.Authenticate(context.Background(), authReq)
	if err != nil {
		t.Error(err)
		return
	}

	config.Verifier.ForceUpdateForTest()

	refreshReq := &proto.RefreshRequest{Sid: sid.Sid, RefreshToken: newRes.RefreshToken}
	refreshRes, err := lHandler.Refresh(context.Background(), refreshReq)
	if err != nil {
		t.Error(err)
		return
	}

	decryptedToken, err := cypher.Decrypt(refreshRes.Sid)
	if err != nil {
		t.Error(err)
		return
	}

	splitedToken := strings.Split(decryptedToken, ".")
	if len(splitedToken) != 2 {
		t.Error("invalid authorized token")
		return
	}
	if len(splitedToken[1]) != 16 {
		t.Error("invalid verification token len")
		return
	}

	_, err = hex.DecodeString(splitedToken[1])
	if err != nil {
		t.Error(err)
		return
	}

	decryptedSid := splitedToken[0]
	splitSid := strings.Split(decryptedSid, "-")
	for _, tokenPart := range splitSid {
		_, err = hex.DecodeString(tokenPart)
		if err != nil {
			t.Error(err)
			return
		}
	}

	verifyReq := &proto.VerifyRequest{
		Sid:              refreshRes.Sid,
		WantedPermission: "c",
	}
	verified, err := lHandler.Verify(context.Background(), verifyReq)
	if err != nil {
		t.Error(err)
		return
	}
	if !verified.Valid {
		t.Error("unable to authorize")
	}
}
