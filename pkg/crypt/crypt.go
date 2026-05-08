package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/config"
)

func Encrypt(text string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(config.AppConfig.Security.Key)
	if err != nil {
		return "", errors.Wrapf(err, "Encrypt failed")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrapf(err, "Encrypt failed")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrapf(err, "Encrypt failed")
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrapf(err, "Encrypt failed")
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(text), nil)
	ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)
	return ciphertextBase64, nil
}

func Decrypt(ciphertextBase64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", errors.Wrapf(err, "Decrypt failed")
	}

	key, err := base64.StdEncoding.DecodeString(config.AppConfig.Security.Key)
	if err != nil {
		return "", errors.Wrapf(err, "Decrypt failed")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrapf(err, "Decrypt failed")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrapf(err, "Decrypt failed")
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.Wrapf(err, "Decrypt failed")
	}

	return string(plaintext), nil
}
