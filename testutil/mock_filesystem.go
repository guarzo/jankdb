package testutil

import (
	"io"
	"os"

	"github.com/guarzo/jankdb"
)

// MockFileSystem is a mock of jankdb.FileSystem.
type MockFileSystem struct {
	ReadFileFunc   func(path string) ([]byte, error)
	OpenFileFunc   func(path string, flag int, perm os.FileMode) (*os.File, error)
	WriteFileFunc  func(path string, data []byte, perm os.FileMode) error
	StatFunc       func(path string) (os.FileInfo, error)
	OpenFunc       func(path string) (io.ReadCloser, error)
	MkdirAllFunc   func(path string, perm os.FileMode) error
	RemoveFunc     func(path string) error
	IsNotExistFunc func(err error) bool
	ReadDirFunc    func(dir string) ([]os.DirEntry, error)
	CreateFunc     func(path string) (*os.File, error)
	RenameFunc     func(src, dst string) error
}

// Compile-time check that MockFileSystem implements jankdb.FileSystem
var _ jankdb.FileSystem = (*MockFileSystem)(nil)

func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	if m.ReadFileFunc == nil {
		return nil, os.ErrNotExist
	}
	return m.ReadFileFunc(path)
}
func (m *MockFileSystem) OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	if m.OpenFileFunc == nil {
		return nil, os.ErrNotExist
	}
	return m.OpenFileFunc(path, flag, perm)
}
func (m *MockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	if m.WriteFileFunc == nil {
		return os.ErrInvalid
	}
	return m.WriteFileFunc(path, data, perm)
}
func (m *MockFileSystem) Stat(path string) (os.FileInfo, error) {
	if m.StatFunc == nil {
		return nil, os.ErrNotExist
	}
	return m.StatFunc(path)
}
func (m *MockFileSystem) Open(path string) (io.ReadCloser, error) {
	if m.OpenFunc == nil {
		return nil, os.ErrNotExist
	}
	return m.OpenFunc(path)
}
func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	if m.MkdirAllFunc == nil {
		return nil
	}
	return m.MkdirAllFunc(path, perm)
}
func (m *MockFileSystem) Remove(path string) error {
	if m.RemoveFunc == nil {
		return nil
	}
	return m.RemoveFunc(path)
}
func (m *MockFileSystem) IsNotExist(err error) bool {
	if m.IsNotExistFunc == nil {
		return os.IsNotExist(err)
	}
	return m.IsNotExistFunc(err)
}
func (m *MockFileSystem) ReadDir(dir string) ([]os.DirEntry, error) {
	if m.ReadDirFunc == nil {
		return nil, nil
	}
	return m.ReadDirFunc(dir)
}
func (m *MockFileSystem) Create(path string) (*os.File, error) {
	if m.CreateFunc == nil {
		return nil, os.ErrNotExist
	}
	return m.CreateFunc(path)
}
func (m *MockFileSystem) Rename(src, dst string) error {
	if m.RenameFunc == nil {
		return nil
	}
	return m.RenameFunc(src, dst)
}
