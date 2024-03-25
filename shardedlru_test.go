// nolint: dupl
package freelru

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// TestShardedRaceCondition tests that the sharded LRU is safe to use concurrently.
// Test with 'go test . -race'.
func TestShardedRaceCondition(t *testing.T) {
	const CAP = 4

	lru, err := NewShardedDefault[uint64, int](CAP)
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

func TestShardedLRUMetrics(t *testing.T) {
	cache, _ := NewShardedDefault[uint64, uint64](1)
	testMetrics(t, cache)

	lru, _ := NewShardedDefault[string, struct{}](7, time.Second*7)
	m := lru.Metrics()
	FatalIf(t, m.Capacity != 7, "Unexpected capacity: %d (!= %d)", m.Capacity, 7)
	FatalIf(t, m.Lifetime != "7s", "Unexpected lifetime: %s (!= %s)", m.Lifetime, "7s")

	lru.ResetMetrics()
	m = lru.Metrics()
	FatalIf(t, m.Capacity != 7, "Unexpected capacity: %d (!= %d)", m.Capacity, 7)
	FatalIf(t, m.Lifetime != "7s", "Unexpected lifetime: %s (!= %s)", m.Lifetime, "7s")
}

func TestStressWithLifetime(t *testing.T) {
	const CAP = 1024

	lru, err := NewShardedDefault[string, int](CAP, time.Millisecond*10)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	const NTHREADS = 10
	const RUNS = 1000

	wg := sync.WaitGroup{}

	for i := 0; i < NTHREADS; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < RUNS; i++ {
				lru.Add(fmt.Sprintf("key-%d", rand.Int()%1000), rand.Int()) //nolint:gosec
				time.Sleep(time.Millisecond * 1)
			}
			wg.Done()
		}()
	}

	for i := 0; i < NTHREADS; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < RUNS; i++ {
				_, _ = lru.Get(fmt.Sprintf("key-%d", rand.Int()%1000)) //nolint:gosec
				time.Sleep(time.Millisecond * 1)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
