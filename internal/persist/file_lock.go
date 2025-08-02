package persist

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// FileLock represents an exclusive lock on the application data directory
type FileLock struct {
	path string
	file *os.File
}

// AcquireLock attempts to acquire an exclusive lock for the application
func AcquireLock(basePath string) (*FileLock, error) {
	lockPath := filepath.Join(basePath, ".lock")
	
	// Check if lock file exists and if the process is still running
	if existingPID, exists := checkExistingLock(lockPath); exists {
		if isProcessRunning(existingPID) {
			return nil, fmt.Errorf("another instance is already running (PID: %d)", existingPID)
		}
		// Stale lock file, remove it
		os.Remove(lockPath)
	}
	
	// Create lock file exclusively
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("failed to acquire lock: another instance may be running")
		}
		return nil, fmt.Errorf("failed to create lock file: %w", err)
	}

	// Write PID and timestamp to lock file
	lockInfo := fmt.Sprintf("%d\n%s\n%s\n", os.Getpid(), runtime.GOOS, time.Now().Format(time.RFC3339))
	if _, err := file.WriteString(lockInfo); err != nil {
		file.Close()
		os.Remove(lockPath)
		return nil, fmt.Errorf("failed to write lock info: %w", err)
	}
	
	if err := file.Sync(); err != nil {
		file.Close()
		os.Remove(lockPath)
		return nil, fmt.Errorf("failed to sync lock file: %w", err)
	}

	return &FileLock{path: lockPath, file: file}, nil
}

// Release releases the lock and removes the lock file
func (fl *FileLock) Release() error {
	if fl.file != nil {
		fl.file.Close()
	}
	return os.Remove(fl.path)
}

// checkExistingLock reads an existing lock file and returns the PID if valid
func checkExistingLock(lockPath string) (int, bool) {
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return 0, false
	}
	
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return 0, false
	}
	
	pid, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, false
	}
	
	return pid, true
}

// isProcessRunning checks if a process with the given PID is running
func isProcessRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	
	switch runtime.GOOS {
	case "windows":
		// On Windows, try to open the process
		proc, err := os.FindProcess(pid)
		if err != nil {
			return false
		}
		// On Windows, FindProcess always succeeds for any PID
		// We need to use a different approach - try to wait with timeout
		done := make(chan struct{})
		go func() {
			proc.Wait()
			close(done)
		}()
		
		select {
		case <-done:
			// Process exited, so it wasn't running
			return false
		case <-time.After(1 * time.Millisecond):
			// Process is still running
			return true
		}
	default:
		// On Unix-like systems, send signal 0 to check if process exists
		// This doesn't actually send a signal, just checks if we can
		proc, err := os.FindProcess(pid)
		if err != nil {
			return false
		}
		err = proc.Signal(syscall.Signal(0))
		return err == nil
	}
}