package app

import (
	"container/list"
	"sync"
)

type CacheItem struct {
    Key   string
    Value interface{}
}

type LRUCache struct {
    capacity int
    items    map[string]*list.Element
    evictList *list.List
    mu       sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        items:    make(map[string]*list.Element),
        evictList: list.New(),
    }
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if element, found := c.items[key]; found {
        c.evictList.MoveToFront(element)
        return element.Value.(*CacheItem).Value, true
    }
    return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if element, found := c.items[key]; found {
        c.evictList.MoveToFront(element)
        element.Value.(*CacheItem).Value = value
        return
    }

    if c.evictList.Len() >= c.capacity {
        c.evict()
    }

    item := &CacheItem{Key: key, Value: value}
    element := c.evictList.PushFront(item)
    c.items[key] = element
}

func (c *LRUCache) evict() {
    element := c.evictList.Back()
    if element != nil {
        c.evictList.Remove(element)
        item := element.Value.(*CacheItem)
        delete(c.items, item.Key)
    }
}