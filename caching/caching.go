package caching

import (
	"time"
)

type DistributedCacheOption struct {
	KeyPrefix              string
	DefaultExpires         time.Duration
	GlobalCacheEntryOption CacheEntryOption
}

type CacheEntryOption struct {
	AbsoluteExpiration              *time.Time
	AbsoluteExpirationRelativeToNow *time.Duration
	SlidingExpiration               *time.Duration
}

type DistributedCacheInterface[T any] interface {
	Get(key string) (T, error)
	Set(key string, in T, option *CacheEntryOption) error
	Delete(key string) error
}

type CacheInterface interface {
	Get(name string, key string) ([]byte, error)
	Set(name string, key string, data []byte, expiration time.Time) error
	Delete(name string, key string) error
}
