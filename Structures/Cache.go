package Structures

import (
	"sync"
	"time"
)

	func CreateCache(defaultLifetime int) Cache {

		myCache := Cache{ DefaultLifetime: defaultLifetime }
		myCache.Setup()
		return myCache
	}

	type Cache struct {

		DefaultLifetime int
		lock sync.RWMutex
		storage map[string]CacheItem
	}

	func (c *Cache) Setup() {

		c.lock = sync.RWMutex{}
		c.storage = make(map[string]CacheItem)
	}

	func (c *Cache) Write(name, content string, expiry int) {

		// Parcel this up
		item := CacheItem{ ID: name, Value: content, Created: time.Now() }

		// Set the expiry for this
		if(expiry < 1) {

			expiry = c.DefaultLifetime
		}

		item.Expiry = time.Now().Add(time.Second * time.Duration(expiry))

		// Store this
		c.lock.Lock()

			c.storage[name] = item

		c.lock.Unlock()
	}

	func (c *Cache) Read(name string) string {

		// Fetch this out of storage
		c.lock.RLock()

			item, exists := c.storage[name]

		c.lock.RUnlock()

		if(!exists) {

			return ""
		}

		// Is this fresh enough?
		if(time.Now().After(item.Expiry)) {

			c.lock.Lock()

				delete(c.storage, name)

			c.lock.Unlock()
			return ""
		}

		return item.Value
	}

	func (c *Cache) ReadAll() []string {

		c.lock.RLock()
		defer c.lock.RUnlock()

		myItems := make([]string, 0)
		for _, item := range c.storage {

			if(time.Now().After(item.Expiry)) {

				continue
			}

			myItems = append(myItems, item.Value)
		}

		return myItems
	}



	type CacheItem struct {

		ID, Value string
		Created time.Time
		Expiry time.Time
	}