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

type Verifier struct {
	previous           *[]byte
	current            *[]byte
	locker             *sync.RWMutex
	startTime          *time.Time
	forceUpdateForTest chan struct{}
	ticker             *time.Ticker
	quit               chan struct{}
}

// This function should only be used for testing purposes.

func NewVerifier() *Verifier {
	startTime := time.Now()
	ticker := time.NewTicker(12 * time.Hour)
	quit := make(chan struct{})

	return &Verifier{
		previous:           CreateRandomBytes(),
		current:            CreateRandomBytes(),
		locker:             &sync.RWMutex{},
		startTime:          &startTime,
		forceUpdateForTest: make(chan struct{}),
		ticker:             ticker,
		quit:               quit,
	}
}

func (v *Verifier) Start() {
	for {
		select {
		//Should be only used for testing purposes.
		case <-v.forceUpdateForTest:
			v.update()
			v.forceUpdateForTest <- struct{}{}
		case <-v.ticker.C:
			v.update()
		case <-v.quit:
			v.ticker.Stop()
			return
		}
	}
}

func (v *Verifier) ForceUpdateForTest() {
	v.forceUpdateForTest <- struct{}{}
	<-v.forceUpdateForTest
}

func (v *Verifier) Verify(pwd string) (err error) {
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

func (v *Verifier) GetExpiration() int64 {
	v.locker.RLock()
	defer v.locker.RUnlock()

	now := time.Now()
	diff := now.Sub(*v.startTime)
	return diff.Milliseconds()
}

func (v *Verifier) GetCurrent() string {
	v.locker.RLock()
	defer v.locker.RUnlock()

	s := hex.EncodeToString(*v.current)
	return s
}

func (v *Verifier) GetPrevious() string {
	v.locker.RLock()
	defer v.locker.RUnlock()

	s := hex.EncodeToString(*v.previous)
	return s
}

func (v *Verifier) update() {
	v.locker.Lock()
	defer v.locker.Unlock()

	*v.previous = *v.current
	*v.current = *CreateRandomBytes()
	*v.startTime = time.Now()
}

func CreateRandomBytes() *[]byte {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	return &randomBytes
}
