package cacher

import (
	"time"
)

type ExampleCache struct {
	Key string
	D   time.Duration
	F   func() error
}

func (n ExampleCache) GetKey() string {
	return n.Key
}

func (n ExampleCache) GetCheckDuration() time.Duration {
	return n.D
}

func (n ExampleCache) GetNewCacheData() error {
	return n.F()
}