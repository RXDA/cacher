package cacher

import (
	"fmt"
	"testing"
	"time"
)

var cache1 string
var cache2 string

func TestRegisterAutoCache(t *testing.T) {
	go func() {
		for {
			// not write lock for brevity
			fmt.Println("t1: ", cache1)
			fmt.Println("t2: ", cache2)
			time.Sleep(time.Second)
		}
	}()

	c := ExampleCache{
		Key: "test1",
		D:   time.Second * 10,
		F: func() error {
			// write cache logic here
			// 若多副本运行，有需要则可添加redis锁
			cache1 = time.Now().String()
			return nil
		},
	}

	c2 := ExampleCache{
		Key: "test2",
		D:   time.Second * 5,
		F: func() error {
			// write cache logic here
			cache2 = time.Now().Format(time.RFC3339Nano)
			return nil
		},
	}
	RegisterAutoCache(c)
	time.Sleep(time.Second * 3)

	RegisterAutoCache(c2)
	time.Sleep(time.Hour)
}
