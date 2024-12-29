package jankdb_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/guarzo/jankdb"
	"github.com/guarzo/jankdb/testutil"
)

func TestStore_Load_NoFile(t *testing.T) {
	mockFS := &testutil.MockFileSystem{
		StatFunc: func(path string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		},
		IsNotExistFunc: func(err error) bool {
			return os.IsNotExist(err)
		},
	}

	s, err := jankdb.NewStore[string](mockFS, "/basepath", jankdb.StoreOptions{
		FileName: "test.json",
	})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	if err = s.Load(); err != nil {
		t.Errorf("expected no error when file is missing, got %v", err)
	}
}

func TestStore_Save_Encrypted(t *testing.T) {
	mockFS := &testutil.MockFileSystem{
		StatFunc: func(path string) (os.FileInfo, error) {
			// pretend file doesn't exist
			return nil, os.ErrNotExist
		},
		IsNotExistFunc: func(err error) bool {
			return os.IsNotExist(err)
		},
		WriteFileFunc: func(path string, data []byte, perm os.FileMode) error {
			// we can check if data is "encrypted"
			if len(data) == 0 {
				return errors.New("no data written")
			}
			return nil
		},
		RenameFunc: func(src, dst string) error {
			return nil
		},
	}

	s, _ := jankdb.NewStore[string](mockFS, "/base", jankdb.StoreOptions{
		FileName:      "securedata.json",
		EncryptionKey: "pass123",
	})

	s.Set("top-secret")
	if err := s.Save(); err != nil {
		t.Errorf("Save failed: %v", err)
	}
}

func TestStore_RoundTrip(t *testing.T) {
	mockFS := &testutil.MockFileSystem{
		// We'll store a buffer in memory to simulate reading/writing
	}

	// A simple in-memory store for the data
	var writtenData []byte

	mockFS.StatFunc = func(path string) (os.FileInfo, error) {
		// if we have writtenData, pretend file exists
		if writtenData != nil {
			return mockFileInfo{name: path}, nil
		}
		return nil, os.ErrNotExist
	}
	mockFS.IsNotExistFunc = func(err error) bool {
		return errors.Is(err, os.ErrNotExist)
	}
	mockFS.WriteFileFunc = func(path string, data []byte, perm os.FileMode) error {
		writtenData = append([]byte(nil), data...)
		return nil
	}
	mockFS.ReadFileFunc = func(path string) ([]byte, error) {
		return writtenData, nil
	}
	mockFS.RenameFunc = func(src, dst string) error { return nil }

	store, _ := jankdb.NewStore[map[string]int](mockFS, "/fakebase", jankdb.StoreOptions{
		FileName: "data.json",
	})
	err := store.Load()
	if err != nil {
		t.Errorf("unexpected load error: %v", err)
	}

	data := store.Get()
	if data == nil {
		data = make(map[string]int)
	}
	data["foo"] = 42
	store.Set(data)

	if err := store.Save(); err != nil {
		t.Errorf("unexpected save error: %v", err)
	}

	// Now re-load
	store2, _ := jankdb.NewStore[map[string]int](mockFS, "/fakebase", jankdb.StoreOptions{
		FileName: "data.json",
	})
	err = store2.Load()
	if err != nil {
		t.Errorf("unexpected load error: %v", err)
	}
	loadedMap := store2.Get()
	if loadedMap["foo"] != 42 {
		t.Errorf("expected 42, got %d", loadedMap["foo"])
	}
}

// mockFileInfo implements os.FileInfo for testing.
type mockFileInfo struct {
	name string
}

func (m mockFileInfo) Name() string           { return m.name }
func (m mockFileInfo) Size() int64            { return 0 }
func (m mockFileInfo) Mode() os.FileMode      { return 0644 }
func (m mockFileInfo) ModTime() (t time.Time) { return }
func (m mockFileInfo) IsDir() bool            { return false }
func (m mockFileInfo) Sys() interface{}       { return nil }
