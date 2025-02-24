package main

import (
	"container/list"
	"fmt"
	"sync"
)

// Определяем структуру для хранения пары ключ-значение.
type entry struct {
	key   string
	value string
}

// CacheLRU хранит данные кэша, список для LRU и мапу для быстрого доступа.
type CacheLRU struct {
	mu        sync.RWMutex
	data      map[string]*list.Element
	listElems *list.List
	cacheSize int
}

// NewCache создает новый кэш с заданным размером.
func NewCache(cacheSize int) *CacheLRU {
	return &CacheLRU{
		data:      make(map[string]*list.Element),
		listElems: list.New(),
		cacheSize: cacheSize,
	}
}

// Set добавляет элемент в кэш или обновляет существующий.
func (c *CacheLRU) Set(k string, v string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если элемент уже существует, обновляем его значение и перемещаем в начало.
	if elem, ok := c.data[k]; ok {
		elem.Value.(*entry).value = v
		c.listElems.MoveToFront(elem)
		return
	}

	// Создаем новый элемент и добавляем его в начало списка.
	elem := c.listElems.PushFront(&entry{key: k, value: v})
	c.data[k] = elem

	// Если кэш переполнен, удаляем наименее недавно использованный элемент.
	if c.listElems.Len() > c.cacheSize {
		lastElem := c.listElems.Back()
		if lastElem != nil {
			c.listElems.Remove(lastElem)
			keyToDelete := lastElem.Value.(*entry).key
			delete(c.data, keyToDelete)
		}
	}
}

// Get возвращает значение по ключу и обновляет его позицию (перемещает в начало).
func (c *CacheLRU) Get(k string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.data[k]; ok {
		// Перемещаем элемент в начало, так как он недавно использовался.
		c.listElems.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return "", false
}

func main() {
	cache := NewCache(2)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if v, ok := cache.Get("key2"); ok {
		fmt.Println("key2:", v)
	}

	// Добавляем третий элемент — кэш переполнен, удалится наименее используемый.
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
