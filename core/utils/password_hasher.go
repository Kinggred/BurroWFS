package utils

import (
	"burrowfs/core/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func EncodePassword(password string) (string, error) {
	block, err := aes.NewCipher(config.CONFIG.HashedSecret)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(password), nil)
	combined := append(nonce, ciphertext...)
	return base64.URLEncoding.EncodeToString(combined), nil
}

func DecodePassword(hash string) (string, error) {
	combined, err := base64.URLEncoding.DecodeString(hash)
	if err != nil {
		return "", err
	}

	if len(combined) < 12 {
		return "", errors.New("hash is too short")
	}
	nonce := combined[:12]
	ciphertext := combined[12:]
	block, err := aes.NewCipher(config.CONFIG.HashedSecret)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
