package main

import (
	"fmt"
	"sync"
)

// CacheLRU реализует LRU-кэш с использованием map и среза.
type CacheLRU struct {
	mu        sync.Mutex
	cacheSize int
	data      map[string]string
	order     []string // порядок ключей: 0-й элемент — самый недавно использованный
}

// NewCache создаёт новый кэш заданного размера.
func NewCache(cacheSize int) *CacheLRU {
	return &CacheLRU{
		cacheSize: cacheSize,
		data:      make(map[string]string),
		order:     make([]string, 0, cacheSize),
	}
}

// moveToFront перемещает ключ в начало среза order.
func (c *CacheLRU) moveToFront(key string) {
	index := -1
	for i, k := range c.order {
		if k == key {
			index = i
			break
		}
	}
	// Если ключ найден, удаляем его из текущей позиции.
	if index != -1 {
		c.order = append(c.order[:index], c.order[index+1:]...)
	}
	// Добавляем ключ в начало.
	c.order = append([]string{key}, c.order...)
}

// Set добавляет новый элемент или обновляет существующий, перемещая его в начало.
func (c *CacheLRU) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если элемент уже существует, обновляем значение и перемещаем ключ в начало.
	if _, exists := c.data[key]; exists {
		c.data[key] = value
		c.moveToFront(key)
		return
	}

	// Если кэш переполнен, удаляем наименее недавно использованный элемент.
	if len(c.order) >= c.cacheSize {
		lruKey := c.order[len(c.order)-1]
		delete(c.data, lruKey)
		c.order = c.order[:len(c.order)-1]
	}

	// Добавляем новый элемент.
	c.data[key] = value
	c.moveToFront(key)
}

// Get возвращает значение по ключу и обновляет его позицию в срезе.
func (c *CacheLRU) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, exists := c.data[key]
	if !exists {
		return "", false
	}
	c.moveToFront(key)
	return value, true
}

func main() {
	cache := NewCache(2)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if v, ok := cache.Get("key2"); ok {
		fmt.Println("key2:", v)
	}

	// Добавляем третий элемент – кэш переполнен, поэтому удалится наименее недавно использованный ("key1").
	cache.Set("key3", "value3")

	if v, ok := cache.Get("key1"); ok {
		fmt.Println("key1:", v)
	} else {
		fmt.Println("key1 not found")
	}

	if v, ok := cache.Get("key3"); ok {
		fmt.Println("key3:", v)
	}
}
