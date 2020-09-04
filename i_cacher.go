package cacher

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
	"time"
)
// TimerCache 定时缓存
type TimerCache interface {
	// 获取缓存key
	GetKey() string
	// 获取检查时间间隔
	GetCheckDuration() time.Duration
	// 获取新数据并保存到缓存
	GetNewCacheData() error
}

var once = sync.Once{}

var startChan = make(chan struct{})

var caches = struct {
	cache []cacheTicker
	*sync.RWMutex
}{
	RWMutex: &sync.RWMutex{},
}

type cacheTicker struct {
	TimerCache
	*time.Ticker
}

func init() {
	once.Do(startAutoCache)
}

func startAutoCache() {
	go func() {
		<-startChan
		for {
			cases := make([]reflect.SelectCase, len(caches.cache))
			for i, ch := range caches.cache {
				cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.C)}
			}
			chosen, _, ok := reflect.Select(cases)
			if ok{
				caches.RLock()
				ch := caches.cache[chosen]
				caches.RUnlock()
				err := ch.GetNewCacheData()
				if err != nil {
					logrus.Errorf("get cache data error, key: %s, error: %s", ch.GetKey(), err.Error())
				}
			}
		}
	}()
}

// RegisterAutoCache return job index
func RegisterAutoCache(c TimerCache) int {
	caches.Lock()
	defer caches.Unlock()
	caches.cache = append(caches.cache, cacheTicker{
		TimerCache: c,
		Ticker:     time.NewTicker(c.GetCheckDuration()),
	})
	if len(caches.cache) == 1{
		startChan<- struct{}{}
	}
	return len(caches.cache) - 1
}
