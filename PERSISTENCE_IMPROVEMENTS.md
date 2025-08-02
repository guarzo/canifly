# Persistence Layer Improvements

## Overview
This document outlines recommended improvements to the current JSON-based persistence layer in CanIFly. These enhancements maintain the simplicity of the current approach while adding robustness and reliability.

## Current Architecture Summary
- JSON files for data storage
- In-memory caching for performance
- Domain-specific stores (Account, Config, EVE data)
- Mutex-protected concurrent access

## Recommended Improvements

### 1. Atomic Write Operations
Prevent data corruption during writes by implementing atomic file operations.

```go
// persist/atomic_writer.go
package persist

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

// AtomicWriteJSON writes data to a file atomically using a temporary file and rename
func AtomicWriteJSON(fs FileSystem, path string, data interface{}) error {
    dir := filepath.Dir(path)
    if err := fs.MkdirAll(dir, os.ModePerm); err != nil {
        return fmt.Errorf("failed to create directories: %w", err)
    }

    // Write to temporary file
    tempPath := path + ".tmp"
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal JSON: %w", err)
    }

    if err := fs.WriteFile(tempPath, jsonData, 0644); err != nil {
        return fmt.Errorf("failed to write temp file: %w", err)
    }

    // Atomic rename
    if err := os.Rename(tempPath, path); err != nil {
        // Cleanup temp file on failure
        fs.Remove(tempPath)
        return fmt.Errorf("failed to rename temp file: %w", err)
    }

    return nil
}
```

### 2. Data Versioning and Migration
Add versioning to handle schema changes gracefully.

```go
// persist/versioned_store.go
package persist

type VersionedData struct {
    Version int         `json:"version"`
    Data    interface{} `json:"data"`
}

type Migration func(oldData []byte, fromVersion int) ([]byte, error)

type VersionedStore struct {
    currentVersion int
    migrations     map[int]Migration
}

func (vs *VersionedStore) Load(fs FileSystem, path string, target interface{}) error {
    data, err := fs.ReadFile(path)
    if err != nil {
        return err
    }

    var versioned VersionedData
    if err := json.Unmarshal(data, &versioned); err != nil {
        // Assume version 1 if no version field
        versioned = VersionedData{Version: 1, Data: data}
    }

    // Apply migrations if needed
    for v := versioned.Version; v < vs.currentVersion; v++ {
        if migration, ok := vs.migrations[v]; ok {
            data, err = migration(data, v)
            if err != nil {
                return fmt.Errorf("migration from v%d failed: %w", v, err)
            }
        }
    }

    return json.Unmarshal(data, target)
}
```

### 3. Write-Ahead Logging (WAL)
Implement a simple WAL for crash recovery.

```go
// persist/wal.go
package persist

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type Operation struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"`      // "write", "delete"
    Path      string    `json:"path"`
    Data      []byte    `json:"data,omitempty"`
    Timestamp time.Time `json:"timestamp"`
}

type WAL struct {
    path   string
    file   *os.File
    writer *bufio.Writer
    fs     FileSystem
}

func NewWAL(fs FileSystem, basePath string) (*WAL, error) {
    walPath := filepath.Join(basePath, "wal.log")
    file, err := os.OpenFile(walPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    return &WAL{
        path:   walPath,
        file:   file,
        writer: bufio.NewWriter(file),
        fs:     fs,
    }, nil
}

func (w *WAL) LogOperation(op Operation) error {
    op.Timestamp = time.Now()
    data, err := json.Marshal(op)
    if err != nil {
        return err
    }

    if _, err := w.writer.Write(data); err != nil {
        return err
    }
    if err := w.writer.WriteByte('\n'); err != nil {
        return err
    }

    return w.writer.Flush()
}

func (w *WAL) Replay() error {
    // Read and replay operations from WAL
    // Implementation details...
    return nil
}

func (w *WAL) Close() error {
    if err := w.writer.Flush(); err != nil {
        return err
    }
    return w.file.Close()
}
```

### 4. File Locking
Prevent multiple process access to the same data files.

```go
// persist/file_lock.go
package persist

import (
    "fmt"
    "os"
    "path/filepath"
)

type FileLock struct {
    path string
    file *os.File
}

func AcquireLock(basePath string) (*FileLock, error) {
    lockPath := filepath.Join(basePath, ".lock")
    file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
    if err != nil {
        if os.IsExist(err) {
            return nil, fmt.Errorf("another instance is already running")
        }
        return nil, err
    }

    // Write PID to lock file
    fmt.Fprintf(file, "%d\n", os.Getpid())
    file.Sync()

    return &FileLock{path: lockPath, file: file}, nil
}

func (fl *FileLock) Release() error {
    if fl.file != nil {
        fl.file.Close()
    }
    return os.Remove(fl.path)
}
```

### 5. Compression for Large Data
Add optional compression for skill data and other large datasets.

```go
// persist/compressed_store.go
package persist

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
    "io"
)

func SaveCompressedJSON(fs FileSystem, path string, data interface{}) error {
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }

    var buf bytes.Buffer
    gz := gzip.NewWriter(&buf)
    if _, err := gz.Write(jsonData); err != nil {
        return err
    }
    if err := gz.Close(); err != nil {
        return err
    }

    return fs.WriteFile(path+".gz", buf.Bytes(), 0644)
}

func LoadCompressedJSON(fs FileSystem, path string, target interface{}) error {
    data, err := fs.ReadFile(path + ".gz")
    if err != nil {
        return err
    }

    gz, err := gzip.NewReader(bytes.NewReader(data))
    if err != nil {
        return err
    }
    defer gz.Close()

    decompressed, err := io.ReadAll(gz)
    if err != nil {
        return err
    }

    return json.Unmarshal(decompressed, target)
}
```

### 6. Integrity Checks
Add checksums to detect data corruption.

```go
// persist/integrity.go
package persist

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

type ChecksummedData struct {
    Data     json.RawMessage `json:"data"`
    Checksum string          `json:"checksum"`
}

func SaveWithChecksum(fs FileSystem, path string, data interface{}) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }

    hash := sha256.Sum256(jsonData)
    checksum := hex.EncodeToString(hash[:])

    wrapped := ChecksummedData{
        Data:     jsonData,
        Checksum: checksum,
    }

    return SaveJsonToFile(fs, path, wrapped)
}

func LoadWithChecksum(fs FileSystem, path string, target interface{}) error {
    var wrapped ChecksummedData
    if err := ReadJsonFromFile(fs, path, &wrapped); err != nil {
        return err
    }

    // Verify checksum
    hash := sha256.Sum256(wrapped.Data)
    checksum := hex.EncodeToString(hash[:])
    if checksum != wrapped.Checksum {
        return fmt.Errorf("checksum mismatch: data may be corrupted")
    }

    return json.Unmarshal(wrapped.Data, target)
}
```

## Implementation Priority

1. **High Priority**
   - Atomic writes (prevents corruption)
   - File locking (prevents concurrent access issues)

2. **Medium Priority**
   - Data versioning (future-proofing)
   - Integrity checks (data validation)

3. **Low Priority**
   - WAL (advanced crash recovery)
   - Compression (optimization for large datasets)

## Migration Strategy

1. Implement atomic writes first - transparent upgrade
2. Add file locking on application startup
3. Gradually introduce versioning for new features
4. Add integrity checks as an optional feature
5. Implement WAL only if crash recovery becomes an issue

## Testing Recommendations

1. **Crash Testing**: Kill process during write operations
2. **Concurrent Access**: Launch multiple instances
3. **Corruption Testing**: Manually corrupt files and test recovery
4. **Performance Testing**: Measure impact of checksums and compression
5. **Migration Testing**: Test version upgrades with real data

## Current Implementation Status

### ✅ Completed
1. **Atomic Write Operations** - Implemented in `atomic_writer.go`
2. **File Locking** - Implemented in `file_lock.go`
3. **Integration** - `SaveJsonToFile` now uses atomic writes transparently
4. **Application Lock** - Added to startup in `cmd/start.go`
5. **Tests** - Comprehensive test coverage for both features

### Current State
- `SaveJsonToFile` is a wrapper that calls `AtomicSaveJsonToFile`
- All existing code benefits from atomic writes without changes
- File locking prevents multiple instances from running

## Migration to Direct Usage Pattern

### Phase 1: Update All Stores to Use AtomicWriteJSON Directly

#### 1.1 Account Store Migration
```go
// internal/persist/account/account_store.go

// OLD PATTERN
func (as *AccountDataStore) saveAccountDataLocked(data model.AccountData) error {
    filePath := filepath.Join(as.basePath, accountFileName)
    if err := persist.SaveJsonToFile(as.fs, filePath, data); err != nil {
        as.logger.WithError(err).Error("Error saving account data")
        return fmt.Errorf("error saving account data: %w", err)
    }
    // ... rest of function
}

// NEW PATTERN - Direct atomic write
func (as *AccountDataStore) saveAccountDataLocked(data model.AccountData) error {
    filePath := filepath.Join(as.basePath, accountFileName)
    if err := persist.AtomicWriteJSON(as.fs, filePath, data); err != nil {
        as.logger.WithError(err).Error("Error saving account data")
        return fmt.Errorf("error saving account data: %w", err)
    }
    // ... rest of function
}
```

#### 1.2 Config Store Migration
```go
// internal/persist/config/config_store.go

// Update all SaveJsonToFile calls to AtomicWriteJSON
// Example:
func (c *ConfigStore) saveConfigDataLocked(configData *model.ConfigData) error {
    filePath := filepath.Join(c.basePath, configFileName)
    // OLD: persist.SaveJsonToFile(c.fs, filePath, configData)
    // NEW:
    if err := persist.AtomicWriteJSON(c.fs, filePath, configData); err != nil {
        c.logger.WithError(err).Error("Error saving config data")
        return err
    }
    // ... rest
}
```

#### 1.3 EVE Store Migration
All stores in `internal/persist/eve/` should be updated similarly.

### Phase 2: Update Non-JSON Writes

#### 2.1 Skill Plan Text Files
```go
// internal/persist/eve/skill_store.go

func (s *SkillStore) SaveSkillPlan(planName string, skills map[string]model.Skill) error {
    // ... build content ...
    
    // NEW: Use atomic write pattern for text files too
    tempPath := planFilePath + ".tmp"
    if err := s.fs.WriteFile(tempPath, []byte(sb.String()), 0644); err != nil {
        return fmt.Errorf("failed to write temp plan file: %w", err)
    }
    
    if err := os.Rename(tempPath, planFilePath); err != nil {
        s.fs.Remove(tempPath)
        return fmt.Errorf("failed to save plan file atomically: %w", err)
    }
    
    // ... rest
}
```

### Phase 3: Remove Backward Compatibility

#### 3.1 Delete Old SaveJsonToFile Function
```go
// internal/persist/fileutils.go

// REMOVE THIS:
func SaveJsonToFile(fs FileSystem, filePath string, source interface{}) error {
    return AtomicSaveJsonToFile(fs, filePath, source)
}

// REMOVE THIS TOO:
func AtomicSaveJsonToFile(fs FileSystem, filePath string, source interface{}) error {
    return AtomicWriteJSON(fs, filePath, source)
}

// KEEP ONLY AtomicWriteJSON as the single source of truth
```

#### 3.2 Update All Import References
Search and replace across the codebase:
- `persist.SaveJsonToFile` → `persist.AtomicWriteJSON`
- `persist.AtomicSaveJsonToFile` → `persist.AtomicWriteJSON`

### Phase 4: Add Atomic Pattern to All File Operations

#### 4.1 Create Generic Atomic Write
```go
// persist/atomic_writer.go

// AtomicWriteFile writes any data atomically (not just JSON)
func AtomicWriteFile(fs FileSystem, path string, data []byte, perm os.FileMode) error {
    dir := filepath.Dir(path)
    if err := fs.MkdirAll(dir, os.ModePerm); err != nil {
        return fmt.Errorf("failed to create directories for %s: %w", path, err)
    }

    tempPath := path + ".tmp"
    if err := fs.WriteFile(tempPath, data, perm); err != nil {
        return fmt.Errorf("failed to write temp file %s: %w", tempPath, err)
    }

    if err := os.Rename(tempPath, path); err != nil {
        fs.Remove(tempPath)
        return fmt.Errorf("failed to rename temp file %s to %s: %w", tempPath, path, err)
    }

    return nil
}
```

### Phase 5: Enhanced Error Handling

#### 5.1 Add Retry Logic for Transient Failures
```go
// persist/atomic_writer.go

func AtomicWriteJSONWithRetry(fs FileSystem, path string, data interface{}, maxRetries int) error {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        if err := AtomicWriteJSON(fs, path, data); err != nil {
            lastErr = err
            // Check if error is retryable
            if isTransientError(err) {
                time.Sleep(time.Millisecond * 100 * time.Duration(i+1))
                continue
            }
            return err
        }
        return nil
    }
    return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

## Implementation Checklist

### Stores to Update
- [ ] `internal/persist/account/account_store.go`
- [ ] `internal/persist/account/login_state_store.go`
- [ ] `internal/persist/config/config_store.go`
- [ ] `internal/persist/config/app_state_store.go`
- [ ] `internal/persist/eve/cache_store.go`
- [ ] `internal/persist/eve/deleted_store.go`
- [ ] `internal/persist/eve/eve_profiles_store.go`
- [ ] `internal/persist/eve/skill_store.go`
- [ ] `internal/persist/eve/system_store.go`

### Code Cleanup
- [ ] Remove `SaveJsonToFile` wrapper function
- [ ] Remove `AtomicSaveJsonToFile` alias
- [ ] Update all imports and function calls
- [ ] Add atomic writes to non-JSON files
- [ ] Update tests to use new patterns

### Documentation Updates
- [ ] Update code comments
- [ ] Update API documentation
- [ ] Add migration notes to CHANGELOG

## Benefits After Migration

1. **Consistency**: Single atomic write pattern across all persistence
2. **Performance**: No double function call overhead
3. **Clarity**: Direct usage makes atomic behavior explicit
4. **Maintainability**: Less code, fewer abstractions
5. **Safety**: All writes are atomic by default

## Rollback Plan

If issues arise during migration:
1. Git revert the specific commits
2. Re-add wrapper functions temporarily
3. Fix issues in isolated branches
4. Re-attempt migration with fixes

## Conclusion

This migration plan moves the codebase from implicit atomic writes (through wrappers) to explicit atomic writes. This makes the code more maintainable and the atomic behavior more obvious to developers.