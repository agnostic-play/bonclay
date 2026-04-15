package common_services

type EncryptionReq struct {
	EncryptionType  string `json:"encryptionType"`  // "AES GCM", "AES CBC", "AES CBC IV", or "AES Imeg"
	Key             string `json:"key"`             // hex, base64, or raw
	KeyDeriveMethod string `json:"keyDeriveMethod"` // "No SUM", "sha256", "sha512", "md5", "sha1", "hmac-sha256"
	RawValue        string `json:"rawValue"`        // plaintext (enc) or ciphertext (dec) - base64/hex/raw
	OutputFormat    string `json:"outputFormat"`    // "hex" or "base64" (encrypt output format)
	InputFormat     string `json:"inputFormat"`     // "hex" or "base64" (encrypt output format)
	IV              string `json:"iv,omitempty"`    // REQUIRED for "AES CBC IV"/"AES Imeg" (hex/base64/raw)
}

type EncryptionResp struct {
	Output any `json:"output"` // encrypted (hex/base64) or decrypted plaintext (base64)
}
