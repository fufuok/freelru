package freelru

// NewDefault constructs an LRU with the given capacity of elements.
// The hash function calculates a hash value from the keys.
func NewDefault[K comparable, V any](capacity uint32) (*LRU[K, V], error) {
	return New[K, V](capacity, MakeHasher[K]())
}

// NewWithSizeDefault constructs an LRU with the given capacity and size.
// The hash function calculates a hash value from the keys.
// A size greater than the capacity increases memory consumption and decreases the CPU consumption
// by reducing the chance of collisions.
// Size must not be lower than the capacity.
func NewWithSizeDefault[K comparable, V any](capacity, size uint32) (*LRU[K, V], error) {
	return NewWithSize[K, V](capacity, size, MakeHasher[K]())
}

// NewSyncedDefault creates a new thread-safe LRU hashmap with the given capacity.
func NewSyncedDefault[K comparable, V any](capacity uint32) (*SyncedLRU[K, V], error) {
	return NewSynced[K, V](capacity, MakeHasher[K]())
}

func NewSyncedWithSizeDefault[K comparable, V any](capacity, size uint32) (*SyncedLRU[K, V], error) {
	return NewSyncedWithSize[K, V](capacity, size, MakeHasher[K]())
}

// NewShardedDefault creates a new thread-safe sharded LRU hashmap with the given capacity.
func NewShardedDefault[K comparable, V any](capacity uint32) (*ShardedLRU[K, V], error) {
	return NewSharded[K, V](capacity, MakeHasher[K]())
}

func NewShardedWithSizeDefault[K comparable, V any](shards, capacity, size uint32) (*ShardedLRU[K, V], error) {
	return NewShardedWithSize[K, V](shards, capacity, size, MakeHasher[K]())
}
