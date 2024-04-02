package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fufuok/freelru"
)

func main() {
	lru, err := freelru.NewSyncedDefault[int64, int64](8192)
	if err != nil {
		panic(err)
	}

	firstKey := int64(9999)
	lru.Add(9999, 111)

	// Simulate concurrent scenarios
	var wg sync.WaitGroup
	wg.Add(8191)
	for i := 1; i < 8192; i++ {
		i := int64(i)
		go func() {
			defer wg.Done()
			lru.Add(i, i)
		}()
	}
	wg.Wait()

	fmt.Println(lru.Get(123))

	lru.SetOnEvict(func(k int64, v int64) {
		fmt.Printf("evicted(firstKey: %d) %d: %d\n", firstKey, k, v) // evicted(firstKey: 9999) 9999: 111
	})

	lru.AddWithLifetime(8888, 8888, 50*time.Millisecond)
	fmt.Println(lru.Get(8888)) // 8888 true
	<-time.After(50 * time.Millisecond)
	fmt.Println(lru.Get(8888)) // 0 false

	fmt.Println(lru.Get(firstKey)) // 0 false

	lru.PrintStats()
	fmt.Printf("Metrics: %+v", lru.Metrics())
}
