package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"strings"
)

func Encode(value any) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte, dest any) error {
	buf := bytes.NewBuffer(data)
	return gob.NewDecoder(buf).Decode(dest)
}

// Encrypt data with AES-GCM
func EncryptData(plaintext, key []byte, encrypted, nonce *[]byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	*nonce = make([]byte, gcm.NonceSize())
	if _, err := rand.Read(*nonce); err != nil {
		return err
	}

	*encrypted = gcm.Seal(nil, *nonce, plaintext, nil)
	return nil
}

func DecryptData(encrypted, nonce, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, encrypted, nil)
}

func HyphenateAndLower(s string) string {
	hyphenated := strings.ReplaceAll(s, " ", "-")
	return strings.ToLower(hyphenated)
}
