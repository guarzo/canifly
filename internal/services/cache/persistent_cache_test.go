package cache_test

import (
	"testing"
	"time"

	cacheSvc "github.com/guarzo/canifly/internal/services/cache"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/services/storage"
	"github.com/guarzo/canifly/internal/testutil"
)

func newTestStorage(t *testing.T) interfaces.StorageService {
	t.Helper()
	return storage.NewStorageService(t.TempDir(), &testutil.MockLogger{})
}

func TestPersistentCache_GetSet_RoundTrip(t *testing.T) {
	c := cacheSvc.NewPersistentCacheService(newTestStorage(t), &testutil.MockLogger{})
	t.Cleanup(c.Shutdown)

	c.Set("foo", []byte("bar"), time.Hour)
	got, ok := c.Get("foo")
	if !ok {
		t.Fatalf("expected key foo to exist")
	}
	if string(got) != "bar" {
		t.Fatalf("expected bar, got %q", got)
	}
}

func TestPersistentCache_Get_Expired(t *testing.T) {
	c := cacheSvc.NewPersistentCacheService(newTestStorage(t), &testutil.MockLogger{})
	t.Cleanup(c.Shutdown)

	c.Set("foo", []byte("bar"), 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if _, ok := c.Get("foo"); ok {
		t.Fatalf("expected expired entry to be missing")
	}
}

func TestPersistentCache_SaveLoad(t *testing.T) {
	store := newTestStorage(t)
	c1 := cacheSvc.NewPersistentCacheService(store, &testutil.MockLogger{})
	t.Cleanup(c1.Shutdown)
	c1.Set("k1", []byte("v1"), time.Hour)
	if err := c1.SaveCache(); err != nil {
		t.Fatalf("SaveCache: %v", err)
	}

	c2 := cacheSvc.NewPersistentCacheService(store, &testutil.MockLogger{})
	t.Cleanup(c2.Shutdown)
	if err := c2.LoadCache(); err != nil {
		t.Fatalf("LoadCache: %v", err)
	}
	got, ok := c2.Get("k1")
	if !ok || string(got) != "v1" {
		t.Fatalf("expected v1, got %q ok=%v", got, ok)
	}
}

func TestPersistentCache_Shutdown_Idempotent(t *testing.T) {
	c := cacheSvc.NewPersistentCacheService(newTestStorage(t), &testutil.MockLogger{})
	c.Shutdown()
	c.Shutdown() // must not panic
}
