package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type verifier struct {
	previous  *[]byte
	current   *[]byte
	locker    *sync.RWMutex
	startTime *time.Time
}

var Verifier *verifier

func SetupVerifier() {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	startTime := time.Now()

	v := &verifier{
		previous:  &randomBytes,
		current:   &randomBytes,
		locker:    &sync.RWMutex{},
		startTime: &startTime,
	}

	setVerifier(v)

	ticker := time.NewTicker(12 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				Verifier.update()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func setVerifier(v *verifier) {
	Verifier = v
}

func (v *verifier) Verify(pwd string) (err error) {
	v.locker.RLock()
	defer v.locker.RUnlock()

	pwdBytes, err := hex.DecodeString(pwd)
	if err != nil {
		return status.Error(codes.Internal, "error parsing hex string")
	}

	if bytes.Compare(pwdBytes, *v.previous) == 0 {
		return status.Error(codes.PermissionDenied, "expired token")
	}

	if bytes.Compare(pwdBytes, *v.current) == 0 {
		return nil
	}

	return status.Error(codes.Unauthenticated, "invalid token")
}

func (v *verifier) GetExpiration() int64 {
	v.locker.RLock()
	defer v.locker.RUnlock()

	now := time.Now()
	diff := now.Sub(*v.startTime)
	return diff.Milliseconds()
}

func (v *verifier) GetCurrent() string {
	v.locker.RLock()
	defer v.locker.RUnlock()

	s := hex.EncodeToString(*v.current)
	return s
}

func (v *verifier) GetPrevious() string {
	v.locker.RLock()
	defer v.locker.RUnlock()

	s := hex.EncodeToString(*v.previous)
	return s
}

func (v *verifier) update() {
	v.locker.Lock()
	defer v.locker.Unlock()

	*v.previous = *v.current
	newPwd := make([]byte, 8)
	rand.Read(newPwd)
	*v.current = newPwd
	*v.startTime = time.Now()
}
