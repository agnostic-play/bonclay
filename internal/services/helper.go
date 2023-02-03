package services

import (
	"crypto/md5"
	"encoding/hex"
)

func encrypt(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func validateMethod(method string) bool {

	return false
}
