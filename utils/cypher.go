package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

type Cypher []byte

var cypher *Cypher

func SetupCypher() {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)

	Cypher := Cypher(randBytes)

	setCypher(&Cypher)
}

func setCypher(Cypher *Cypher) {
	cypher = Cypher
}

func GetCypher() Cypher {
	return *cypher
}

func (cy *Cypher) Encrypt(s string) (string, error) {
	c, err := aes.NewCipher(*cy)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cypher := gcm.Seal(nonce, nonce, []byte(s), nil)
	value, err := hex.EncodeToString(cypher), nil
	if err != nil {
		return "", err
	}

	return value, nil
}

func (cy *Cypher) Decrypt(s string) (string, error) {
	sBytes, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(*cy)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(sBytes) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := sBytes[:nonceSize], sBytes[nonceSize:]
	decryptedBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}
