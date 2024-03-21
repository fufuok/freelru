// nolint: dupl
package freelru

import (
	"sync"
	"testing"
	"time"
)

// TestSyncedRaceCondition tests that the synced LRU is safe to use concurrently.
// Test with 'go test . -race'.
func TestSyncedRaceCondition(t *testing.T) {
	const CAP = 4

	lru, err := NewSyncedDefault[uint64, int](CAP)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	wg := sync.WaitGroup{}

	call := func(fn func()) {
		wg.Add(1)
		go func() {
			fn()
			wg.Done()
		}()
	}

	call(func() { lru.SetLifetime(0) })
	call(func() { lru.SetOnEvict(nil) })
	call(func() { _ = lru.Len() })
	call(func() { _ = lru.AddWithLifetime(1, 1, 0) })
	call(func() { _ = lru.Add(1, 1) })
	call(func() { _, _ = lru.Get(1) })
	call(func() { _, _ = lru.Peek(1) })
	call(func() { _ = lru.Contains(1) })
	call(func() { _ = lru.Remove(1) })
	call(func() { _ = lru.Keys() })
	call(func() { lru.Purge() })
	call(func() { lru.Metrics() })
	call(func() { _ = lru.ResetMetrics() })
	call(func() { lru.dump() })
	call(func() { lru.PrintStats() })

	wg.Wait()
}

func TestSyncedLRUMetrics(t *testing.T) {
	cache, _ := NewSyncedDefault[uint64, uint64](1)
	testMetrics(t, cache)

	lru, _ := NewSyncedDefault[string, struct{}](3, time.Second*3)
	m := lru.Metrics()
	FatalIf(t, m.Capacity != 3, "Unexpected capacity: %d (!= %d)", m.Capacity, 1)
	FatalIf(t, m.Lifetime != "3s", "Unexpected lifetime: %s (!= %s)", m.Lifetime, "3s")
}
