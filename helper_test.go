package freelru

import (
	"testing"
)

func TestDefaultLRUAddGet(t *testing.T) {
	const CAP = 1_000_000

	cache, _ := NewDefault[int, int](CAP * 2)
	m := make(map[int]int, CAP)

	for i := 0; i < CAP; i++ {
		x := int(runtimeFastrand())
		cache.Add(i, x)
		m[i] = x
	}

	cache.PrintStats()

	for i := 0; i < CAP; i++ {
		v, ok := cache.Get(i)
		if !ok {
			t.Fatalf("key: %d not found", i)
		}
		if v != m[i] {
			t.Fatalf("key: %d, value: %d != %d", i, v, m[i])
		}
	}
}

func TestSyncedDefaultLRUAddGet(t *testing.T) {
	const CAP = 1_000_000

	cache, _ := NewSyncedDefault[int, int](CAP * 2)
	m := make(map[int]int, CAP)

	for i := 0; i < CAP; i++ {
		x := int(runtimeFastrand())
		cache.Add(i, x)
		m[i] = x
	}

	cache.PrintStats()

	for i := 0; i < CAP; i++ {
		v, ok := cache.Get(i)
		if !ok {
			t.Fatalf("key: %d not found", i)
		}
		if v != m[i] {
			t.Fatalf("key: %d, value: %d != %d", i, v, m[i])
		}
	}
}

func TestShardedDefaultLRUAddGet(t *testing.T) {
	const CAP = 1_000_000

	cache, _ := NewShardedDefault[int, int](CAP * 2)
	m := make(map[int]int, CAP)

	for i := 0; i < CAP; i++ {
		x := int(runtimeFastrand())
		cache.Add(i, x)
		m[i] = x
	}

	cache.PrintStats()

	for i := 0; i < CAP; i++ {
		v, ok := cache.Get(i)
		if !ok {
			t.Fatalf("key: %d not found", i)
		}
		if v != m[i] {
			t.Fatalf("key: %d, value: %d != %d", i, v, m[i])
		}
	}
}
