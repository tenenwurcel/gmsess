package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

type verifier struct {
	previous *[]byte
	current  *[]byte
	locker   *sync.RWMutex
}

var Verifier *verifier

func SetupVerifier() {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	v := &verifier{
		previous: &randomBytes,
		current:  &randomBytes,
		locker:   &sync.RWMutex{},
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
		return
	}

	if bytes.Compare(pwdBytes, *v.previous) == 0 {
		return errors.New("Expired token!")
	}

	if bytes.Compare(pwdBytes, *v.current) == 0 {
		return nil
	}

	return errors.New("Invalid token!")
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
	fmt.Println(hex.EncodeToString(*v.current))
}
