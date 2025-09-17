package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type EncryptionToolsEntity struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}

type ToolsService interface {
	EncryptConfig(ctx context.Context, req EncryptionToolsEntity) (string, error)
	DecryptConfig(ctx context.Context, req EncryptionToolsEntity) (string, error)
}

func (cont serviceContainer) EncryptConfig(ctx context.Context, req EncryptionToolsEntity) (string, error) {
	decodeString, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return "", err
	}

	key := sha256.Sum256([]byte(req.Key))
	encryptedConf, err := encrypt(key[:], decodeString)
	if err != nil {
		return "", err
	}

	return encryptedConf, nil
}

func (cont serviceContainer) DecryptConfig(ctx context.Context, req EncryptionToolsEntity) (string, error) {
	key := sha256.Sum256([]byte(req.Key))
	conf, err := decrypt(key[:], req.Data)
	if err != nil || conf == nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(conf), nil
}

func encrypt(key, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-GCM.
func decrypt(key []byte, cipherTextBase64 string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
