package freelru

import (
	"time"
)

// NewDefault constructs an LRU with the given capacity of elements.
// The hash function calculates a hash value from the keys.
func NewDefault[K comparable, V any](capacity uint32, lifetime ...time.Duration) (*LRU[K, V], error) {
	return NewWithSizeDefault[K, V](capacity, capacity, lifetime...)
}

// NewWithSizeDefault constructs an LRU with the given capacity and size.
// The hash function calculates a hash value from the keys.
// A size greater than the capacity increases memory consumption and decreases the CPU consumption
// by reducing the chance of collisions.
// Size must not be lower than the capacity.
func NewWithSizeDefault[K comparable, V any](capacity, size uint32, lifetime ...time.Duration) (*LRU[K, V], error) {
	c, err := NewWithSize[K, V](capacity, size, MakeHasher[K]())
	if err != nil {
		return nil, err
	}
	if len(lifetime) > 0 && lifetime[0] > 0 {
		c.SetLifetime(lifetime[0])
	}
	return c, err
}

// NewSyncedDefault creates a new thread-safe LRU hashmap with the given capacity.
func NewSyncedDefault[K comparable, V any](capacity uint32, lifetime ...time.Duration) (*SyncedLRU[K, V], error) {
	return NewSyncedWithSizeDefault[K, V](capacity, capacity, lifetime...)
}

func NewSyncedWithSizeDefault[K comparable, V any](capacity, size uint32, lifetime ...time.Duration) (
	*SyncedLRU[K, V], error,
) {
	c, err := NewSyncedWithSize[K, V](capacity, size, MakeHasher[K]())
	if err != nil {
		return nil, err
	}
	if len(lifetime) > 0 && lifetime[0] > 0 {
		c.SetLifetime(lifetime[0])
	}
	return c, err
}

// NewShardedDefault creates a new thread-safe sharded LRU hashmap with the given capacity.
func NewShardedDefault[K comparable, V any](capacity uint32, lifetime ...time.Duration) (*ShardedLRU[K, V], error) {
	c, err := NewSharded[K, V](capacity, MakeHasher[K]())
	if err != nil {
		return nil, err
	}
	if len(lifetime) > 0 && lifetime[0] > 0 {
		c.SetLifetime(lifetime[0])
	}
	return c, err
}

func NewShardedWithSizeDefault[K comparable, V any](shards, capacity, size uint32, lifetime ...time.Duration) (
	*ShardedLRU[K, V], error,
) {
	c, err := NewShardedWithSize[K, V](shards, capacity, size, MakeHasher[K]())
	if err != nil {
		return nil, err
	}
	if len(lifetime) > 0 && lifetime[0] > 0 {
		c.SetLifetime(lifetime[0])
	}
	return c, err
}
