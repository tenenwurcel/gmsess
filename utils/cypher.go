package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Cypher []byte

func (cy *Cypher) Encrypt(s string) (string, error) {
	c, err := aes.NewCipher(*cy)
	if err != nil {
		return "", status.Error(codes.Internal, "error creating aes cypher")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", status.Error(codes.Internal, "error creating aes cypher")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", status.Error(codes.Internal, "error reading nonce bytes")
	}

	cypher := gcm.Seal(nonce, nonce, []byte(s), nil)
	value, err := hex.EncodeToString(cypher), nil
	if err != nil {
		return "", status.Error(codes.Internal, "error reading bytes to string")
	}

	return value, nil
}

func (cy *Cypher) Decrypt(s string) (string, error) {
	sBytes, err := hex.DecodeString(s)
	if err != nil {
		return "", status.Error(codes.Internal, "error decoding hexadecimal token.")
	}

	c, err := aes.NewCipher(*cy)
	if err != nil {
		return "", status.Error(codes.Internal, "error creating aes cypher")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", status.Error(codes.Internal, "error creating gcm cypher")
	}

	nonceSize := gcm.NonceSize()
	if len(sBytes) < nonceSize {
		return "", status.Error(codes.Internal, "ciphertext too short")
	}

	nonce, ciphertext := sBytes[:nonceSize], sBytes[nonceSize:]
	decryptedBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}

	return string(decryptedBytes), nil
}
