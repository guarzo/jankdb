package jankdb

import (
	"io"
	"os"
)

// FileSystem is your abstraction over the file system.
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	OpenFile(path string, flag int, perm os.FileMode) (*os.File, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	Stat(path string) (os.FileInfo, error)
	Open(path string) (io.ReadCloser, error)
	MkdirAll(path string, perm os.FileMode) error
	Remove(path string) error
	IsNotExist(err error) bool
	ReadDir(dir string) ([]os.DirEntry, error)
	Create(dst string) (*os.File, error)
	Rename(src, dst string) error
}

// OSFileSystem is a real implementation that calls the `os` package.
type OSFileSystem struct{}

func (OSFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
func (OSFileSystem) OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, flag, perm)
}
func (OSFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}
func (OSFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
func (OSFileSystem) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
func (OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
func (OSFileSystem) Remove(path string) error {
	return os.Remove(path)
}
func (OSFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
func (OSFileSystem) ReadDir(dir string) ([]os.DirEntry, error) {
	return os.ReadDir(dir)
}
func (OSFileSystem) Create(dir string) (*os.File, error) {
	return os.Create(dir)
}
func (OSFileSystem) Rename(src, dst string) error {
	return os.Rename(src, dst)
}
