package jankdb

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

// Store[T] is a generic store for any type T.
type Store[T any] struct {
	mu sync.RWMutex

	fs       FileSystem
	basePath string

	// If you want a subdirectory (like "loot", "app", etc.)
	subDir   string
	fileName string

	data T

	cache *Cache[T]

	// Backup old file as .bak before overwriting
	enableBackup bool

	// If non-empty => encrypt on write, decrypt on read
	encryptionKey string
}

// StoreOptions defines the parameters for customizing a Store.
type StoreOptions struct {
	SubDir       string
	FileName     string
	EnableBackup bool

	UseCache          bool
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration

	// If not empty, we do AES-GCM encryption using this passphrase
	EncryptionKey string
}

// NewStore creates a new Store[T].
func NewStore[T any](fs FileSystem, basePath string, opts StoreOptions) (*Store[T], error) {
	s := &Store[T]{
		fs:            fs,
		basePath:      basePath,
		subDir:        opts.SubDir,
		fileName:      opts.FileName,
		enableBackup:  opts.EnableBackup,
		encryptionKey: opts.EncryptionKey,
	}

	if opts.UseCache {
		s.cache = NewCache[T](opts.DefaultExpiration, opts.CleanupInterval)
	}

	return s, nil
}

// Load reads T from the file. If `encryptionKey` is set, we decrypt the file first.
func (s *Store[T]) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.filePath()
	if _, err := s.fs.Stat(path); s.fs.IsNotExist(err) {
		// No file => do nothing
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Read raw bytes
	bytes, err := s.fs.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var tmp T
	if s.encryptionKey == "" {
		// Plain JSON
		if err := json.Unmarshal(bytes, &tmp); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	} else {
		// Decrypt
		plaintext, err := DecryptData(s.encryptionKey, string(bytes))
		if err != nil {
			return fmt.Errorf("failed to decrypt data: %w", err)
		}
		if err := json.Unmarshal(plaintext, &tmp); err != nil {
			return fmt.Errorf("failed to unmarshal decrypted JSON: %w", err)
		}
	}

	s.data = tmp
	if s.cache != nil {
		s.cache.Set("all", s.data)
	}
	return nil
}

// Save writes T to disk, using atomic write & optional .bak backup.
// If encryptionKey is not empty, data is encrypted before writing.
func (s *Store[T]) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.filePath()
	dir := filepath.Dir(path)

	// Ensure directory
	if err := s.fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// 1) Marshal data to JSON
	bytes, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data as JSON: %w", err)
	}

	if s.encryptionKey != "" {
		// 2) Encrypt the JSON
		encrypted, err := EncryptData(s.encryptionKey, bytes)
		if err != nil {
			return fmt.Errorf("failed to encrypt data: %w", err)
		}
		// 3) Perform atomic write with the encrypted data
		if err := atomicWriteFile(s.fs, path, []byte(encrypted), s.enableBackup); err != nil {
			return fmt.Errorf("failed to write encrypted data: %w", err)
		}
	} else {
		// 2) Plain JSON => atomic write
		if err := atomicWriteFile(s.fs, path, bytes, s.enableBackup); err != nil {
			return fmt.Errorf("failed to write JSON data: %w", err)
		}
	}

	return nil
}

// Get returns the in-memory data.
func (s *Store[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data
}

// Set replaces the entire in-memory data.
func (s *Store[T]) Set(val T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = val
	if s.cache != nil {
		s.cache.Set("all", val)
	}
}

// filePath -> /basePath/subDir/fileName
func (s *Store[T]) filePath() string {
	if s.subDir == "" {
		return filepath.Join(s.basePath, s.fileName)
	}
	return filepath.Join(s.basePath, s.subDir, s.fileName)
}

// atomicWriteFile writes data to .tmp, optionally backups .bak, then renames to final.
func atomicWriteFile(fs FileSystem, finalPath string, data []byte, enableBackup bool) error {
	tmpPath := finalPath + ".tmp"
	bakPath := finalPath + ".bak"

	// Backup old file
	if enableBackup {
		if _, err := fs.Stat(finalPath); err == nil {
			if err := fs.Rename(finalPath, bakPath); err != nil {
				return fmt.Errorf("failed to rename old file to .bak: %w", err)
			}
		}
	}

	// Write to .tmp
	if err := fs.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write .tmp file: %w", err)
	}

	// Rename tmp -> final
	if err := fs.Rename(tmpPath, finalPath); err != nil {
		return fmt.Errorf("failed to rename %s -> %s: %w", tmpPath, finalPath, err)
	}

	return nil
}
