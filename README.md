# jankdb

[![Build & Test CI](https://github.com/guarzo/jankdb/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/guarzo/jankdb/actions/workflows/ci.yml)
[![Release Workflow](https://github.com/guarzo/jankdb/actions/workflows/release.yml/badge.svg)](https://github.com/guarzo/jankdb/actions/workflows/release.yml)


A **lightweight** and **extensible** key-value store for Go, supporting:
1. **Atomic Writes** (writes to a temporary file, then renames to final).
2. **Optional JSON Encryption** (AES-GCM, with a scrypt-derived key).
3. **Optional .bak backups** before overwriting existing files.
4. **Optional In-Memory Caching** (via [patrickmn/go-cache](https://github.com/patrickmn/go-cache)).

> **Repository**: [github.com/guarzo/jankdb](https://github.com/guarzo/jankdb)

---

## Features

- **Simple Go API** – Use `jankdb.Store[T]` to store any Go type (`map`, `struct`, `[]Something`, etc.).
- **Atomic File Writes** – Prevent data corruption by writing to a `.tmp` file and renaming once complete.
- **Automatic Backups** – Optionally rename the old file to `.bak` before overwriting.
- **Encryption at Rest** – Enable encryption by specifying an `EncryptionKey`; data is transparently encrypted/decrypted.
- **In-Memory Caching** – Speed up reads with a configurable TTL cache.
- **File System Abstraction** – Built-in `OSFileSystem` for real I/O, or provide your own mock `FileSystem` for testing.

---

## Installation

```bash
go get github.com/guarzo/jankdb
```

Then import it in your code:

```go
import "github.com/guarzo/jankdb"
```

---

## Getting Started

### 1) Plain JSON Storage

If you just want to store your data as plain JSON (no encryption, no cache):

```go
package main

import (
"fmt"

    "github.com/guarzo/jankdb"
)

type LootSplit struct {
// ... your fields
}

func main() {
fs := jankdb.OSFileSystem{}
opts := jankdb.StoreOptions{
SubDir:       "loot",
FileName:     "loot_split.json",
EnableBackup: true,  // create a .bak file before overwriting
UseCache:     false, // no in-memory caching
EncryptionKey: "",   // empty => no encryption
}

    lootStore, err := jankdb.NewStore[[]LootSplit](fs, "/path/to/base", opts)
    if err != nil {
        panic(err)
    }

    // Load existing data from disk if present
    if err := lootStore.Load(); err != nil {
        panic(err)
    }

    // Access or modify
    splits := lootStore.Get()
    fmt.Println("Current loot splits:", splits)

    // Save changes
    if err := lootStore.Save(); err != nil {
        panic(err)
    }
}
```

**Key Steps**:

1. Instantiate a `Store[T]` with the desired `FileSystem`, base path, and `StoreOptions`.
2. Call `Load()` once to read existing data from disk.
3. Call `Get()` and `Set(...)` to read and modify the in-memory data.
4. Call `Save()` to write changes back to disk atomically.

---

### 2) Encrypted JSON Storage

By providing an `EncryptionKey`, you enable built-in **AES-GCM** encryption. Data is encrypted before writing to disk and decrypted upon load.

```go
type Identity struct {
MainID string
// ...
}

func main() {
fs := jankdb.OSFileSystem{}
opts := jankdb.StoreOptions{
SubDir:         "app",
FileName:       "identities.json.enc",
EnableBackup:   true,
UseCache:       true,                      // also enable caching
DefaultExpiration: 30 * time.Minute,       // cache TTL
CleanupInterval:   5 * time.Minute,
EncryptionKey:   "MySuperSecretPassword",  // set a passphrase
}

    store, _ := jankdb.NewStore[map[string]Identity](fs, "/secure/path", opts)
    _ = store.Load() // decrypts if file exists

    data := store.Get()
    if data == nil {
        data = make(map[string]Identity)
    }
    data["12345"] = Identity{MainID: "SomeID"}
    store.Set(data)

    _ = store.Save() // automatically encrypts & atomically writes to file
}
```

**Notes**:
- **Do not** commit your `EncryptionKey` to source control.
- This built-in approach uses a **scrypt**-derived AES-GCM scheme. For production-grade security, review your key management, scrypt parameters, and consider using more advanced cryptographic solutions.

---

### 3) Caching

`jankdb` can store a copy of your data in memory for faster reads. You set:

- **`UseCache: true`** to enable caching.
- **`DefaultExpiration`** to define how long an item remains in cache.
- **`CleanupInterval`** to define how often expired items are purged.

```go
opts := jankdb.StoreOptions{
FileName:          "mydata.json",
UseCache:          true,
DefaultExpiration: 15 * time.Minute,
CleanupInterval:   5 * time.Minute,
// ...
}
myStore, _ := jankdb.NewStore[MyData](fs, "/some/dir", opts)
_ = myStore.Load()   // populates cache
cachedVal := myStore.Get()  // immediate
```

> **Warning**: If multiple processes modify the same file, you’ll need your own synchronization to keep caches consistent.

---

## Project Status

`jankdb` is **experimental** and was created to reduce boilerplate for simple data persistence. It is not intended to replace heavyweight databases. Use it for small to medium “configuration” or “state” files, especially where portability and simplicity matter more than scale.

---

## Potential Use Cases

- **Configuration** or **User Preferences** in a CLI or server application.
- **Encrypted secrets** for small-scale projects that don’t require a dedicated secrets manager.
- **Caching** ephemeral data for improved performance.

---

## Limitations and Security

- **Single-writer model**: `jankdb` is designed for a single process or thread writing to the file at a time.
- **Encryption**: This module only provides basic AES-GCM encryption with scrypt-based key derivation. In high-security contexts, you may need more rigorous key management and encryption strategies.
- **Concurrency**: The `Store[T]` type is guarded by a `sync.RWMutex`, so concurrent reads and writes from multiple goroutines should work, but the underlying data type `T` itself must be safe to manipulate from multiple threads (or you must carefully manage concurrent updates).
- **Performance**: Each `Save()` operation rewrites the entire file. If your data is very large, you may need a different approach (e.g., partial updates, a real database).

---

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Contact

**Author**: [@guarzo](https://github.com/guarzo)  
**Project Link**: [github.com/guarzo/jankdb](https://github.com/guarzo/jankdb)
