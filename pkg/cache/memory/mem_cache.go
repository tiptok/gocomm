// inspired from https://github.com/patrickmn/go-cache
package cache

import (
	"fmt"
	. "github.com/tiptok/gocomm/pkg/cache/model"
	"sync"
	"time"
)

type MemCache struct {
	*memcache
}

type memcache struct {
	items   sync.Map
	janitor *janitor
}

// NewMemCache memcache will scan all objects per clean interval, and delete expired key.
func NewMemCache(ci time.Duration) *MemCache {
	c := &memcache{
		items: sync.Map{},
	}
	C := &MemCache{c}
	//if ci > 0 {
	//	runJanitor(c, ci)
	//	runtime.SetFinalizer(C, stopJanitor)
	//}
	return C
}

// Get an item from the memcache. Returns the item or nil, and a bool indicating whether the key was found.
func (c *memcache) Get(k string, obj interface{}) (*Item, error) {
	tmp, found := c.items.Load(k)
	if !found {
		return nil, nil
	}
	item := tmp.(*Item)
	//obj = item.Object
	//fmt.Println("memcache get cache item:",item.String())
	return item, nil
}

func (c *memcache) Set(k string, it *Item) error {
	if it == nil {
		return fmt.Errorf("memcahce:set cache item is nil")
	}
	c.items.Store(k, it)
	//fmt.Println("memcache set cache item:",it.String())
	return nil
}

// Delete an item from the memcache. Does nothing if the key is not in the memcache.
func (c *memcache) Delete(key string) error {
	c.delete(key)
	//fmt.Println("memcache delete cache item:",k)
	return nil
}
func (c *memcache) delete(key string) (interface{}, bool) {
	c.items.Delete(key)
	return nil, false
}

// Delete all expired items from the memcache.
func (c *memcache) DeleteExpired() {
	c.items.Range(func(key, value interface{}) bool {
		v := value.(*Item)
		k := key.(string)
		// delete outdate for memory cahce
		if v.Expire() {
			c.delete(k)
		}
		return true
	})
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *memcache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(c *MemCache) {
	c.janitor.stop <- true
}

func runJanitor(c *memcache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}
