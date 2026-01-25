package common_services

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type EncryptionService interface {
	EncryptConfig(ctx context.Context, req EncryptionReq) (EncryptionResp, error)
	DecryptConfig(ctx context.Context, req EncryptionReq) (EncryptionResp, error)
}

type encryptionService struct{}

func NewEncryptionService() EncryptionService {
	return &encryptionService{}
}

func (s *encryptionService) EncryptConfig(ctx context.Context, req EncryptionReq) (EncryptionResp, error) {
	// Decode plaintext (try base64; fallback to raw string)
	plain, err := base64.StdEncoding.DecodeString(req.RawValue)
	if err != nil {
		plain = []byte(req.RawValue)
	}

	key, err := buildAESKeyFromInput(req.Key, req.KeyDeriveMethod)
	if err != nil {
		return EncryptionResp{}, err
	}

	var out []byte

	mode := strings.ToLower(req.EncryptionType)
	switch mode {
	case "aesgcm":
		out, err = encryptAESGCM(key, plain)

	case "aescbc":
		// Legacy: random IV generated and prefixed to ciphertext
		out, err = encryptAESCBC(key, plain)

	case "aescbciv":
		// Explicit IV CBC (no IV prefix in output)
		if strings.TrimSpace(req.IV) == "" {
			return EncryptionResp{}, errors.New("iv is required for AES CBC IV")
		}
		out, err = encryptAESCBCWithIV(key, plain, []byte(req.IV))

	case "aesimeg":
		// EXACT Imeg encryption: returns base64-encoded ciphertext (no IV prefix)
		if strings.TrimSpace(req.IV) == "" {
			return EncryptionResp{}, errors.New("iv is required for AES Imeg")
		}

		if lenIV := len(req.IV); lenIV != 16 {
			return EncryptionResp{}, errors.New("iv is need more than 16 bytes")
		}

		if lenKey := len(req.Key); lenKey != 16 {
			return EncryptionResp{}, errors.New("iv is need more than 16 bytes")
		}

		// call EncryptImeg EXACTLY as given
		val, err2 := EncryptImeg(plain, key, []byte(req.IV))
		if err2 != nil {
			return EncryptionResp{}, err2
		}

		out = val

	default:
		err = fmt.Errorf("unsupported encryptionType: %s", req.EncryptionType)
	}

	if err != nil {
		return EncryptionResp{}, err
	}

	return EncryptionResp{Output: encodeOutput(out, req.OutputFormat)}, nil
}

func (s *encryptionService) DecryptConfig(ctx context.Context, req EncryptionReq) (resp EncryptionResp, err error) {
	key, err := buildAESKeyFromInput(req.Key, req.KeyDeriveMethod)
	if err != nil {
		return EncryptionResp{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("decrypt panic: %v", r)
			resp = EncryptionResp{}
		}
	}()

	// Parse input ciphertext
	var ciphertext []byte
	ciphertext, err = decodeInputDecrypter(req.RawValue, req.InputFormat)
	if err != nil {
		return EncryptionResp{}, fmt.Errorf("invalid rawValue: %w", err)
	}

	mode := strings.ToLower(req.EncryptionType)
	var plain []byte
	switch mode {
	case "aesgcm":
		plain, err = decryptAESGCM(key, ciphertext)

	case "aescbc":
		// Legacy: IV prefixed to ciphertext
		plain, err = decryptAESCBC(key, ciphertext)

	case "aescbciv":
		// Explicit IV CBC (no IV in ciphertext)
		if strings.TrimSpace(req.IV) == "" {
			return EncryptionResp{}, errors.New("iv is required for AES CBC IV")
		}
		plain, err = decryptAESCBCWithIV(key, ciphertext, []byte(req.IV))

	case "aesimeg":
		// Imeg variant needs IV supplied separately
		if strings.TrimSpace(req.IV) == "" {
			return EncryptionResp{}, errors.New("iv is required for AES Imeg")
		}
		plain, err = DecryptImeg(ciphertext, key, []byte(req.IV))

	default:
		err = fmt.Errorf("unsupported encryptionType: %s", req.EncryptionType)
	}

	if err != nil {
		return EncryptionResp{}, err
	}

	// Return plaintext as base64 (safe for arbitrary binary)
	resp = EncryptionResp{Output: base64.StdEncoding.EncodeToString(plain)}
	return
}

// Build AES key from input key + method. Supports arbitrary input length by sizing to 16/24/32.
func buildAESKeyFromInput(keyStr, method string) ([]byte, error) {
	material := []byte(keyStr)
	switch strings.ToLower(strings.TrimSpace(method)) {
	case "no sum", "nosum", "raw", "":
		return []byte(material), nil
	case "sha256":
		sum := sha256.Sum256(material)
		return sum[:], nil
	case "sha512":
		sum := sha512.Sum512(material)
		return sum[:], nil
	case "md5":
		sum := md5.Sum(material)
		return sum[:], nil
	case "sha1":
		sum := sha1.Sum(material)
		return sum[:], nil
	case "hmac-sha256":
		// Static salt for compatibility; replace with your own secret management.
		mac := hmac.New(sha256.New, []byte("secret-hmac-salt"))
		mac.Write(material)
		return mac.Sum(nil), nil
	default:
		return nil, fmt.Errorf("unsupported key_derive_method: %s", method)
	}
}

func encryptAESGCM(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	// Prefix nonce
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptAESGCM(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ns := gcm.NonceSize()
	if len(ciphertext) < ns {
		return nil, errors.New("ciphertext too short")
	}
	nonce, data := ciphertext[:ns], ciphertext[ns:]
	return gcm.Open(nil, nonce, data, nil)
}

// ============================
// AES-CBC (PKCS7)
// ============================

// Legacy behavior: random IV generated and prefixed to ciphertext.
func encryptAESCBC(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bs := block.BlockSize()
	plaintext = pkcs7Pad(plaintext, bs)

	iv := make([]byte, bs)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	ct := make([]byte, len(plaintext))
	mode.CryptBlocks(ct, plaintext)

	return append(iv, ct...), nil
}

// Legacy behavior: IV is read from the prefix of the ciphertext.
func decryptAESCBC(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(ciphertext) < bs {
		return nil, errors.New("ciphertext too short")
	}
	iv, data := ciphertext[:bs], ciphertext[bs:]
	if len(data)%bs != 0 {
		return nil, errors.New("ciphertext length must be a multiple of block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	pt := make([]byte, len(data))
	mode.CryptBlocks(pt, data)
	return pkcs7Unpad(pt)
}

// Explicit-IV CBC (no IV prefix).
func encryptAESCBCWithIV(key, plaintext, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(iv) != bs {
		return nil, fmt.Errorf("iv must be %d bytes", bs)
	}
	plaintext = pkcs7Pad(plaintext, bs)

	mode := cipher.NewCBCEncrypter(block, iv)
	ct := make([]byte, len(plaintext))
	mode.CryptBlocks(ct, plaintext)
	return ct, nil
}

func decryptAESCBCWithIV(key, ciphertext, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(iv) != bs {
		return nil, fmt.Errorf("iv must be %d bytes", bs)
	}
	if len(ciphertext)%bs != 0 {
		return nil, errors.New("ciphertext length must be a multiple of block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	pt := make([]byte, len(ciphertext))
	mode.CryptBlocks(pt, ciphertext)
	return pkcs7Unpad(pt)
}

// ============================
// EXACT Imeg block (UNMODIFIED)
// ============================

func EncryptImeg(plaintext []byte, key, iv []byte) (encoded []byte, err error) {
	bPlaintext := pkcs5Padding(plaintext, aes.BlockSize, len(string(plaintext)))

	cipherData, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	log.Println(iv)
	if len(iv) != cipherData.BlockSize() {
		return nil, fmt.Errorf("%v %v", len(iv), cipherData.BlockSize())
	}

	cipherText := make([]byte, len(bPlaintext))
	stream := cipher.NewCBCEncrypter(cipherData, iv)
	stream.CryptBlocks(cipherText, bPlaintext)

	return cipherText, err
}

func DecryptImeg(cipherText []byte, key, iv []byte) (result []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("%v %v", len(iv), cipherText)
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
	if len(b) == 0 {
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

func pkcs7Unpad(data []byte) ([]byte, error) {
	l := len(data)
	if l == 0 {
		return nil, errors.New("invalid padding size")
	}
	pad := int(data[l-1])
	if pad == 0 || pad > l {
		return nil, errors.New("invalid padding bytes")
	}
	// Verify all padding bytes
	for i := l - pad; i < l; i++ {
		if int(data[i]) != pad {
			return nil, errors.New("invalid padding bytes")
		}
	}
	return data[:l-pad], nil
}

func decodeInputDecrypter(val, inputFormat string) ([]byte, error) {
	if strings.ToLower(inputFormat) == "base64" {
		decodedCal, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, err
		}

		return decodedCal, nil
	}

	decodedCal, err := hex.DecodeString(val)
	if err != nil {
		return nil, err
	}

	// As a last resort, treat as raw bytes
	return decodedCal, nil
}

// encodeOutput: "hex" or "base64" (default base64)
func encodeOutput(outputBody []byte, fmtOpt string) string {
	switch strings.ToLower(strings.TrimSpace(fmtOpt)) {
	case "hex":
		return hex.EncodeToString(outputBody)
	case "base64":
		return base64.StdEncoding.EncodeToString(outputBody)
	default:
		return string(outputBody)
	}
}
