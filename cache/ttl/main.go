package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Profile struct {
	UUID   string
	Name   string
	Orders []*Order
}

type Order struct {
	UUID      string
	Value     any
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CacheProfile struct {
	profile      Profile
	lastUserTime time.Time
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]CacheProfile
	ttl  time.Duration
}

func NewCache(ctx context.Context, ttl time.Duration, initSize int) *Cache {
	cache := &Cache{
		data: make(map[string]CacheProfile, initSize),
		ttl:  ttl,
	}

	ticker := time.NewTicker(ttl)
	// Очистка по ttl и закрытие горутины по контексту
	go func(c *Cache) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for key, cacheProfile := range c.data {
					if time.Now().Sub(cacheProfile.lastUserTime) >= ttl {
						c.mu.Lock()
						delete(c.data, key)
						c.mu.Unlock()
					}
				}
			}
		}
	}(cache)

	return cache
}

func (c *Cache) Set(key string, value Profile) {
	c.mu.Lock()
	c.data[key] = CacheProfile{
		profile:      value,
		lastUserTime: time.Now(),
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (Profile, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	if !ok {
		fmt.Printf("Not found key %s\n", key)

		return Profile{}, false
	}
	c.data[key] = CacheProfile{
		profile:      value.profile,
		lastUserTime: time.Now(),
	}

	return value.profile, true
}

func main() {
	createProfile := func(profileUUID, orderUUID, name string) Profile {
		return Profile{
			UUID: profileUUID,
			Name: name,
			Orders: []*Order{{
				UUID: orderUUID,
			}},
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cache := NewCache(ctx, time.Second, 10)

	// empty cache
	_, ok := cache.Get("uuid_1")
	if ok {
		log.Fatal("TEST empty cache - FAIL")
	} else {
		fmt.Println("TEST empty cache - OK")
	}

	// set and get
	profile1 := createProfile("123", "345", "Alex")
	cache.Set(profile1.UUID, profile1)

	profile2 := createProfile("333", "555", "Frank")
	cache.Set(profile2.UUID, profile2)

	cachedProfile1, ok := cache.Get(profile1.UUID)
	if !ok {
		log.Fatal("TEST set and get: FAIL")
	}

	cachedProfile2, ok := cache.Get(profile2.UUID)
	if !ok {
		log.Fatal("TEST set and get: FAIL")
	}
	if profile1.Orders[0] == cachedProfile1.Orders[0] &&
		profile2.Orders[0] == cachedProfile2.Orders[0] {
		fmt.Println("TEST set and get: OK")
	} else {
		fmt.Println("TEST set and get: FAIL")
	}

	// set, wait by ttl and get
	profile1 = createProfile("555", "666", "George")
	cache.Set(profile1.UUID, profile1)

	profile2 = createProfile("999", "001", "Fedor")
	cache.Set(profile2.UUID, profile2)

	_, ok = cache.Get(profile2.UUID)
	if !ok {
		log.Fatal("TEST set, wait by ttl and get: FAIL")
	}

	time.Sleep(time.Millisecond * 700)
	_, ok = cache.Get(profile2.UUID)
	if !ok {
		log.Fatal("TEST set, wait by ttl and get: FAIL")
	}
	time.Sleep(time.Millisecond * 400)

	cachedProfile2, ok = cache.Get(profile2.UUID)
	if !ok {
		log.Fatal("TEST set, wait by ttl and get: FAIL")
	}

	cachedProfile1, ok = cache.Get(profile1.UUID)
	if ok {
		log.Fatal("TEST set, wait by ttl and get: FAIL")
	}
	if profile2.Orders[0] == cachedProfile2.Orders[0] &&
		cachedProfile1.Orders == nil {
		fmt.Println("TEST set, wait by ttl and get: OK")
	} else {
		fmt.Println("TEST set, wait by ttl and get: FAIL")
	}
}
