package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var PreviewDir string

// LRUCache реализация LRU кэша.
type LRUCache struct {
	capacity int
	cache    map[string]*node
	head     *node
	tail     *node
	mutex    sync.Mutex
}

// node структура для двусвязного списка.
type node struct {
	key   string
	value []byte
	prev  *node
	next  *node
}

// NewLRUCache создаёт новый LRU кэш.
func NewLRUCache(capacity int) *LRUCache {
	if PreviewDir == "" {
		PreviewDir = "./previews"
	}

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*node),
	}
}

func (c *LRUCache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if n, ok := c.cache[key]; ok {
		c.moveToFront(n)

		return n.value, true
	}

	return nil, false
}

func (c *LRUCache) Set(key string, value []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if n, ok := c.cache[key]; ok {
		n.value = value
		c.moveToFront(n)

		return nil
	}

	if len(c.cache) >= c.capacity {
		if err := c.removeLRU(); err != nil {
			return fmt.Errorf("error removing preview from disk: %w", err)
		}
	}

	n := &node{
		key:   key,
		value: value,
	}

	if err := c.saveToDisk(key, value); err != nil {
		return fmt.Errorf("error saving preview to disk: %w", err)
	}

	c.cache[key] = n
	c.addToFront(n)

	return nil
}

func (c *LRUCache) moveToFront(n *node) {
	if c.head == n {
		return
	}

	c.remove(n)
	c.addToFront(n)
}

func (c *LRUCache) addToFront(n *node) {
	n.next = c.head
	n.prev = nil

	if c.head != nil {
		c.head.prev = n
	}

	c.head = n

	if c.tail == nil {
		c.tail = n
	}
}

func (c *LRUCache) remove(n *node) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}
}

func (c *LRUCache) removeLRU() error {
	if c.tail != nil {
		if err := c.deleteFromDisk(c.tail.key); err != nil {
			return err
		}

		delete(c.cache, c.tail.key)
		c.remove(c.tail)

		return nil
	}

	return nil
}

func (c *LRUCache) saveToDisk(key string, data []byte) error {
	if err := os.MkdirAll(PreviewDir, 0o750); err != nil {
		return fmt.Errorf("error creating preview directory: %w", err)
	}

	path := filepath.Join(PreviewDir, key)

	return os.WriteFile(path, data, 0o600)
}

func (c *LRUCache) deleteFromDisk(key string) error {
	path := filepath.Join(PreviewDir, key)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error deleting file from disk key:%s, err:%w", key, err)
	}

	return nil
}
