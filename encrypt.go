package jankdb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/scrypt"
	"io"
)

// EncryptData encrypts `plaintext` using a passphrase. Returns a base64-encoded ciphertext.
func EncryptData(passphrase string, plaintext []byte) (string, error) {
	// 1. Generate salt
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to read random salt: %w", err)
	}

	// 2. Derive key from passphrase + salt using scrypt
	key, err := scrypt.Key([]byte(passphrase), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", fmt.Errorf("failed to derive key: %w", err)
	}

	// 3. Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 4. Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to read random nonce: %w", err)
	}

	// 5. Seal
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil) // no additional data

	// 6. Final payload = salt || nonce || ciphertext
	payload := append(salt, nonce...)
	payload = append(payload, ciphertext...)

	// 7. Return base64-encoded
	return base64.StdEncoding.EncodeToString(payload), nil
}

// DecryptData decrypts a base64-encoded ciphertext string using the passphrase.
// Returns the plaintext.
func DecryptData(passphrase, base64CipherText string) ([]byte, error) {
	// 1. Decode base64
	payload, err := base64.StdEncoding.DecodeString(base64CipherText)
	if err != nil {
		return nil, fmt.Errorf("failed to base64-decode ciphertext: %w", err)
	}
	if len(payload) < 16 {
		return nil, fmt.Errorf("payload too short to contain salt")
	}

	// 2. Extract salt
	salt := payload[:16]
	rest := payload[16:]

	// Derive key from pass + salt
	key, err := scrypt.Key([]byte(passphrase), salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// 3. Create AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(rest) < nonceSize {
		return nil, fmt.Errorf("ciphertext missing nonce")
	}

	nonce := rest[:nonceSize]
	ciphertext := rest[nonceSize:]

	// 4. Open (decrypt)
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}
