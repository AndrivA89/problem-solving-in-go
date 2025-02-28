package main

// List представляет узел связного списка
type List struct {
	next  *List
	key   string
	value string
}

// CacheMap представляет хэш-таблицу с цепочечной адресацией
type CacheMap struct {
	arr []*List
}

// NewCacheMap создает новую хэш-таблицу заданного размера
func NewCacheMap(size int) *CacheMap {
	cacheMap := &CacheMap{
		arr: make([]*List, size),
	}
	return cacheMap
}

// getHash вычисляет индекс в массиве для заданного ключа
func (c *CacheMap) getHash(key string) int {
	return len(key) % len(c.arr)
}

// Set добавляет новую пару ключ-значение в хэш-таблицу
func (c *CacheMap) Set(key, value string) {
	index := c.getHash(key)
	if c.arr[index] == nil {
		c.arr[index] = &List{key: key, value: value}
		return
	}

	// Решаем коллизию
	for node := c.arr[index]; node != nil; node = node.next {
		if node.key == key {
			node.value = value // Обновляем значение, если ключ уже есть
			return
		}
		if node.next == nil {
			node.next = &List{key: key, value: value}
			return
		}
	}
}

// Get ищет значение по ключу в хэш-таблице
func (c *CacheMap) Get(key string) (string, bool) {
	index := c.getHash(key)
	for node := c.arr[index]; node != nil; node = node.next {
		if node.key == key {
			return node.value, true
		}
	}
	return "", false
}

// Delete удаляет ключ из хэш-таблицы
func (c *CacheMap) Delete(key string) bool {
	index := c.getHash(key)
	if c.arr[index] == nil {
		return false
	}

	// Если ключ первый в цепочке
	if c.arr[index].key == key {
		c.arr[index] = c.arr[index].next
		return true
	}

	// Поиск и удаление в цепочке
	for node := c.arr[index]; node.next != nil; node = node.next {
		if node.next.key == key {
			node.next = node.next.next
			return true
		}
	}
	return false
}
