package jankdb_test

import (
	"errors"
	"os"
	"testing"

	"github.com/guarzo/jankdb"
	"github.com/guarzo/jankdb/testutil"
)

func TestMockFileSystem(t *testing.T) {
	mockFS := &testutil.MockFileSystem{
		ReadFileFunc: func(path string) ([]byte, error) {
			if path == "valid.txt" {
				return []byte("Hello"), nil
			}
			return nil, os.ErrNotExist
		},
	}

	data, err := mockFS.ReadFile("valid.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", string(data))
	}

	_, err = mockFS.ReadFile("missing.txt")
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestOSFileSystem_Stat(t *testing.T) {
	fs := jankdb.OSFileSystem{}
	_, err := fs.Stat("some-nonexistent-file.txt")
	// Typically, we expect an error
	if err == nil {
		t.Error("expected file to not exist, but got no error")
	}
}
