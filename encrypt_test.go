package jankdb_test

import (
	"strings"
	"testing"

	"github.com/guarzo/jankdb"
)

func TestEncryptDecryptData(t *testing.T) {
	pass := "testpass"
	plain := []byte("secret message")

	encrypted, err := jankdb.EncryptData(pass, plain)
	if err != nil {
		t.Fatalf("EncryptData failed: %v", err)
	}

	// should be base64, not containing plaintext
	if strings.Contains(encrypted, "secret message") {
		t.Error("plaintext found in ciphertext string")
	}

	decrypted, err := jankdb.DecryptData(pass, encrypted)
	if err != nil {
		t.Fatalf("DecryptData failed: %v", err)
	}
	if string(decrypted) != string(plain) {
		t.Errorf("expected %s, got %s", plain, decrypted)
	}
}

func TestDecryptData_BadPassphrase(t *testing.T) {
	pass := "correctpass"
	plain := []byte("mydata")

	encrypted, _ := jankdb.EncryptData(pass, plain)
	_, err := jankdb.DecryptData("wrongpass", encrypted)
	if err == nil {
		t.Error("expected error when decrypting with wrong passphrase, got nil")
	}
}
