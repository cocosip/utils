package cache2go

import (
	"fmt"
	"github.com/muesli/cache2go"
	"time"
)

type Cache struct{}

func (c *Cache) Get(name string, key string) ([]byte, error) {
	table := cache2go.Cache(name)
	item, err := table.Value(key)
	if err != nil {
		return nil, err
	}

	if b, ok := (item.Data()).([]byte); ok {
		return b, nil
	}

	return nil, fmt.Errorf("cache2go convert to []byte fail")
}

func (c *Cache) Set(name string, key string, data []byte, expiration time.Duration) error {
	table := cache2go.Cache(name)
	table.Add(key, expiration, data)
	return nil
}

func (c *Cache) Delete(name string, key string) error {
	table := cache2go.Cache(name)
	if _, err := table.Delete(key); err != nil {
		return err
	}
	return nil
}
