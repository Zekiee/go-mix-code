package lruConcurrentMap

import (
	"container/list"
	"sync"
	"time"
)

type LRUCache struct {
	mutex     sync.Mutex
	maxCount  int
	cacheList *list.List
	cacheMap  map[interface{}]*list.Element
	checkTime time.Duration
	stopChan  chan struct{}
}

type Pair struct {
	key        interface{}
	value      interface{}
	expiredAt  time.Time
	updateTime time.Time
}

func NewLRUCache(capacity int, checkTime time.Duration) *LRUCache {
	c := &LRUCache{
		maxCount:  capacity,
		cacheList: list.New(),
		cacheMap:  make(map[interface{}]*list.Element),
		checkTime: checkTime,
		stopChan:  make(chan struct{}),
	}
	go c.start()
	return c
}

func (c *LRUCache) start() {
	ticker := time.NewTicker(c.checkTime)
	for {
		select {
		case <-ticker.C:
			c.checkExpired()
		case <-c.stopChan:
			ticker.Stop()
			return
		}
	}
}

func (c *LRUCache) Stop() {
	close(c.stopChan)
}

func (c *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, exists := c.cacheMap[key]; exists {
		pair := elem.Value.(*Pair)
		pair.updateTime = time.Now()
		c.cacheList.MoveToFront(elem)
		return pair.value, true
	}

	return nil, false
}

func (c *LRUCache) Set(key interface{}, value interface{}, expiredAt time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, exists := c.cacheMap[key]; exists {
		pair := elem.Value.(*Pair)
		pair.value = value
		pair.expiredAt = expiredAt
		pair.updateTime = time.Now()
		c.cacheList.MoveToFront(elem)
		return
	}

	pair := &Pair{
		key:        key,
		value:      value,
		expiredAt:  expiredAt,
		updateTime: time.Now(),
	}
	elem := c.cacheList.PushFront(pair)
	c.cacheMap[key] = elem

	if c.maxCount != 0 && c.cacheList.Len() > c.maxCount {
		elem := c.cacheList.Back()
		c.removeElement(elem)
	}
}

func (c *LRUCache) removeElement(elem *list.Element) {
	pair := elem.Value.(*Pair)
	delete(c.cacheMap, pair.key)
	c.cacheList.Remove(elem)
}

func (c *LRUCache) checkExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for {
		elem := c.cacheList.Back()
		if elem == nil {
			break
		}

		pair := elem.Value.(*Pair)
		if !pair.expiredAt.IsZero() && time.Now().After(pair.expiredAt) {
			c.removeElement(elem)
		} else {
			break
		}
	}
}
