package services

import (
	`bytes`
	`crypto/aes`
	`crypto/cipher`
	`encoding/base64`
	`errors`
)

func EncryptImeg(plaintext []byte, key, iv []byte) (encoded string, err error) {
	bPlaintext := pkcs5Padding(plaintext, aes.BlockSize, len(string(plaintext)))

	c, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	cipherText := make([]byte, len(bPlaintext))

	// if _, err = io.ReadFull(rand.Reader, iv); err != nil {
	// 	return
	// }

	stream := cipher.NewCBCEncrypter(c, iv)
	stream.CryptBlocks(cipherText, bPlaintext)
	return base64.StdEncoding.EncodeToString(cipherText), err
}

func DecryptImeg(ciphertext string, key, iv []byte) (result []byte, err error) {
	// Remove base64 encoding:
	cipherText, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return
	}

	// Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return
	}

	// Decrypt the message
	stream := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(cipherText))
	stream.CryptBlocks(decrypted, cipherText)

	result = pkcs5Trimming(decrypted)

	return
}

func pkcs5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7Pad(b []byte, blocksize int) []byte {
	if blocksize <= 0 {
		return nil
	}
	if b == nil || len(b) == 0 {
		return nil
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
