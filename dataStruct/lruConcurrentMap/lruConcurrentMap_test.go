package lruConcurrentMap

import (
	"testing"
	"time"
)

func TestNewLRUCache(t *testing.T) {

	cache := NewLRUCache(3, 4*time.Second)

	cache.Set("a", 1, time.Now().Add(time.Second))
	cache.Set("b", 2, time.Now().Add(3*time.Second))
	cache.Set("c", 3, time.Now().Add(10*time.Second))

	if _, ok := cache.Get("a"); !ok {
		println("there is a bug, a")
	}

	if _, ok := cache.Get("b"); !ok {
		println("there is a bug, b")
	}

	if _, ok := cache.Get("c"); !ok {
		println("there is a bug, c")
	}

	time.Sleep(5 * time.Second)

	if _, ok := cache.Get("a"); ok {
		println("there is a bug, a")
	}

	if _, ok := cache.Get("b"); ok {
		println("there is a bug, b")
	}

	if _, ok := cache.Get("c"); !ok {
		println("there is a bug, c")
	}

	cache.Stop()
}
