package eve_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/eve"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCacheStore_EmptyInitially(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := eve.NewCacheStore(logger, fs, basePath)

	// No items initially
	val, found := store.Get("key")
	assert.False(t, found)
	assert.Nil(t, val)

	// Loading from non-existent file should not error
	err := store.LoadApiCache()
	assert.NoError(t, err)
}

func TestCacheStore_SetGet(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := eve.NewCacheStore(logger, fs, basePath)

	data := []byte("testdata")
	store.Set("mykey", data, 10*time.Minute)

	val, found := store.Get("mykey")
	assert.True(t, found)
	assert.Equal(t, data, val)

	val2, found2 := store.Get("otherkey")
	assert.False(t, found2)
	assert.Nil(t, val2)
}

func TestCacheStore_SaveAndLoad(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := eve.NewCacheStore(logger, fs, basePath)

	// Add some items
	store.Set("key1", []byte("value1"), 10*time.Minute)
	store.Set("key2", []byte("value2"), 5*time.Minute)

	err := store.SaveApiCache()
	assert.NoError(t, err)

	// Create a new store and load the data
	newStore := eve.NewCacheStore(logger, fs, basePath)
	err = newStore.LoadApiCache()
	assert.NoError(t, err)

	val, found := newStore.Get("key1")
	assert.True(t, found)
	assert.Equal(t, []byte("value1"), val)

	val2, found2 := newStore.Get("key2")
	assert.True(t, found2)
	assert.Equal(t, []byte("value2"), val2)
}

func TestCacheStore_ExpiredItemsNotLoaded(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := eve.NewCacheStore(logger, fs, basePath)

	// Set an item that expires soon
	store.Set("soonToExpire", []byte("expiring"), 1*time.Millisecond)

	// Wait for it to expire
	time.Sleep(10 * time.Millisecond)

	// Save to file (it should still be in the cache, but expired)
	err := store.SaveApiCache()
	assert.NoError(t, err)

	// Create a new store and load the data
	newStore := eve.NewCacheStore(logger, fs, basePath)
	err = newStore.LoadApiCache()
	assert.NoError(t, err)

	// The expired item should not be set
	val, found := newStore.Get("soonToExpire")
	assert.False(t, found)
	assert.Nil(t, val)
}

func TestCacheStore_LoadInvalidData(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := eve.NewCacheStore(logger, fs, basePath)

	// Write invalid JSON to the cache file
	filePath := filepath.Join(basePath, "cache.json")
	err := os.WriteFile(filePath, []byte("not valid json"), 0600)
	assert.NoError(t, err)

	// Loading should fail
	err = store.LoadApiCache()
	assert.Error(t, err)
}

func TestCacheStore_SaveNoByteValues(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := eve.NewCacheStore(logger, fs, basePath)

	// Insert a non-byte value directly (we can't do that easily via store.Set since it expects []byte)
	// Instead, we'll use the underlying cache for this test scenario:
	store.Set("notBytes", []byte("original"), 10*time.Minute)

	// Modify item in the underlying cache for testing:
	items := store.CacheItemsForTest()
	items["notBytes"] = items["notBytes"] // This is just a no-op; If you need to force a non-byte, you'd have to reflect or break type safety.

	// Actually, to test the "not bytes" scenario realistically:
	// Let's store a non-byte value via store.cache directly.
	store.CacheSetForTest("notBytes2", "stringValue")

	// Now save to file and ensure warning is logged
	// "notBytes2" should be skipped
	err := store.SaveApiCache()
	assert.NoError(t, err)

	// Load again
	newStore := eve.NewCacheStore(logger, fs, basePath)
	err = newStore.LoadApiCache()
	assert.NoError(t, err)

	// "notBytes2" should not be found
	val, found := newStore.Get("notBytes2")
	assert.False(t, found)
	assert.Nil(t, val)

	// "notBytes" was bytes, so it should be found
	val, found = newStore.Get("notBytes")
	assert.True(t, found)
	assert.Equal(t, []byte("original"), val)
}
