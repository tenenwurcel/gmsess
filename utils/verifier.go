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
	previous           *[]byte
	current            *[]byte
	locker             *sync.RWMutex
	startTime          *time.Time
	forceUpdateForTest chan struct{}
}

var Verifier *verifier

func SetupVerifier() {
	startTime := time.Now()

	v := &verifier{
		previous:           CreateRandomBytes(),
		current:            CreateRandomBytes(),
		locker:             &sync.RWMutex{},
		startTime:          &startTime,
		forceUpdateForTest: make(chan struct{}),
	}

	setVerifier(v)

	ticker := time.NewTicker(12 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			//Should be only used for testing purposes.
			case <-Verifier.forceUpdateForTest:
				Verifier.update()
				v.forceUpdateForTest <- struct{}{}
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

func CreateRandomBytes() *[]byte {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	return &randomBytes
}

// This function should only be used for testing purposes.
func (v *verifier) ForceUpdateForTest() {
	v.forceUpdateForTest <- struct{}{}
	<-v.forceUpdateForTest
}

func (v *verifier) Verify(pwd string) (err error) {
	v.locker.RLock()
	defer v.locker.RUnlock()

	pwdBytes, err := hex.DecodeString(pwd)
	if err != nil {
		return status.Error(codes.Internal, "error parsing hex string")
	}

	if bytes.Compare(pwdBytes, *v.current) == 0 {
		return nil
	}

	if bytes.Compare(pwdBytes, *v.previous) == 0 {
		return status.Error(codes.PermissionDenied, "expired token")
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
	*v.current = *CreateRandomBytes()
	*v.startTime = time.Now()
}
