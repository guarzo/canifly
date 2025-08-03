package persist

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestAcquireLock(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lock_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test acquiring lock
	lock, err := AcquireLock(tempDir)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	defer lock.Release()

	// Verify lock file exists
	lockPath := filepath.Join(tempDir, ".lock")
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Fatal("Lock file was not created")
	}

	// Read lock file and verify content
	data, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("Failed to read lock file: %v", err)
	}

	content := string(data)
	if content == "" {
		t.Fatal("Lock file is empty")
	}

	// Verify PID is present
	lines := splitLines(content)
	if len(lines) < 1 {
		t.Fatal("Lock file missing PID")
	}

	pid, err := strconv.Atoi(lines[0])
	if err != nil {
		t.Fatalf("Invalid PID in lock file: %v", err)
	}

	if pid != os.Getpid() {
		t.Fatalf("Wrong PID in lock file: got %d, want %d", pid, os.Getpid())
	}
}

func TestAcquireLockExclusive(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lock_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// First lock
	lock1, err := AcquireLock(tempDir)
	if err != nil {
		t.Fatalf("Failed to acquire first lock: %v", err)
	}
	defer lock1.Release()

	// Try to acquire second lock - should fail
	lock2, err := AcquireLock(tempDir)
	if err == nil {
		lock2.Release()
		t.Fatal("Expected error when acquiring second lock, got nil")
	}

	if err.Error() != "another instance is already running (PID: "+strconv.Itoa(os.Getpid())+")" {
		t.Fatalf("Unexpected error message: %v", err)
	}
}

func TestReleaseLock(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lock_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	lock, err := AcquireLock(tempDir)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	// Release the lock
	err = lock.Release()
	if err != nil {
		t.Fatalf("Failed to release lock: %v", err)
	}

	// Verify lock file is removed
	lockPath := filepath.Join(tempDir, ".lock")
	if _, err := os.Stat(lockPath); err == nil {
		t.Fatal("Lock file still exists after release")
	}

	// Should be able to acquire lock again
	lock2, err := AcquireLock(tempDir)
	if err != nil {
		t.Fatalf("Failed to re-acquire lock after release: %v", err)
	}
	lock2.Release()
}

func TestStaleLockCleanup(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Stale lock cleanup test is flaky on Windows")
	}

	tempDir, err := os.MkdirTemp("", "lock_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a fake lock file with non-existent PID
	lockPath := filepath.Join(tempDir, ".lock")
	stalePID := "99999999\n" + runtime.GOOS + "\n" + time.Now().Format(time.RFC3339)
	if err := os.WriteFile(lockPath, []byte(stalePID), 0644); err != nil {
		t.Fatalf("Failed to create stale lock file: %v", err)
	}

	// Should be able to acquire lock (stale lock should be cleaned up)
	lock, err := AcquireLock(tempDir)
	if err != nil {
		t.Fatalf("Failed to acquire lock with stale lock file: %v", err)
	}
	defer lock.Release()
}

func TestCheckExistingLock(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lock_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	lockPath := filepath.Join(tempDir, ".lock")

	// Test with no lock file
	pid, exists := checkExistingLock(lockPath)
	if exists {
		t.Fatal("checkExistingLock returned true for non-existent file")
	}
	if pid != 0 {
		t.Fatalf("Expected PID 0 for non-existent file, got %d", pid)
	}

	// Test with valid lock file
	validContent := "12345\n" + runtime.GOOS + "\n" + time.Now().Format(time.RFC3339)
	if err := os.WriteFile(lockPath, []byte(validContent), 0644); err != nil {
		t.Fatalf("Failed to create lock file: %v", err)
	}

	pid, exists = checkExistingLock(lockPath)
	if !exists {
		t.Fatal("checkExistingLock returned false for existing file")
	}
	if pid != 12345 {
		t.Fatalf("Expected PID 12345, got %d", pid)
	}

	// Test with invalid content
	if err := os.WriteFile(lockPath, []byte("not-a-number\n"), 0644); err != nil {
		t.Fatalf("Failed to write invalid lock file: %v", err)
	}

	pid, exists = checkExistingLock(lockPath)
	if exists {
		t.Fatal("checkExistingLock returned true for invalid content")
	}
}

func TestIsProcessRunning(t *testing.T) {
	// Test with current process (should be running)
	if !isProcessRunning(os.Getpid()) {
		t.Fatal("isProcessRunning returned false for current process")
	}

	// Test with invalid PID
	if isProcessRunning(-1) {
		t.Fatal("isProcessRunning returned true for negative PID")
	}

	if isProcessRunning(0) {
		t.Fatal("isProcessRunning returned true for PID 0")
	}

	// Test with very high PID (unlikely to exist)
	if isProcessRunning(99999999) {
		t.Fatal("isProcessRunning returned true for non-existent PID")
	}
}

// Helper function to split string into lines (handles different line endings)
func splitLines(s string) []string {
	var lines []string
	line := ""
	for _, r := range s {
		if r == '\n' || r == '\r' {
			if line != "" {
				lines = append(lines, line)
				line = ""
			}
		} else {
			line += string(r)
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}
