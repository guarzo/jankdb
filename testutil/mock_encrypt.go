package testutil

import "fmt"

// FakeEncryptData returns the same data with a prefix "ENCRYPTED:" for demonstration.
func FakeEncryptData(passphrase string, plaintext []byte) (string, error) {
	// For testing, just pretend to encrypt:
	return fmt.Sprintf("ENCRYPTED:%s", string(plaintext)), nil
}

// FakeDecryptData strips off "ENCRYPTED:" prefix.
func FakeDecryptData(passphrase, base64CipherText string) ([]byte, error) {
	const prefix = "ENCRYPTED:"
	if len(base64CipherText) < len(prefix) {
		return []byte(base64CipherText), nil
	}
	if base64CipherText[:len(prefix)] == prefix {
		return []byte(base64CipherText[len(prefix):]), nil
	}
	// Not prefixed => return as is
	return []byte(base64CipherText), nil
}
