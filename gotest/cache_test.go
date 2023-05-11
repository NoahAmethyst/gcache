package gotest

import (
	"github.com/NoahAmethyst/gcache"
	"strconv"
	"sync"
	"testing"
)

func Test_ConcurrentPut(t *testing.T) {
	max := 5000
	localCache := gcache.NewCache[int](max, 10)

	var wait sync.WaitGroup

	for i := 0; i < max*2; i++ {
		wait.Add(1)
		go func(k int, v string) {
			defer wait.Done()
			localCache.Put(k, v)
		}(i, strconv.Itoa(i))
	}

	wait.Wait()

	t.Logf("data size:%d", len(localCache.Keys()))
}

func Test_BatchPut(t *testing.T) {
	max := 5000
	localCache := gcache.NewCache[int](max, 10)

	for i := 0; i < max*2; i++ {
		localCache.Put(i, strconv.Itoa(i))
	}

	keys := localCache.Keys()

	t.Logf("data size:%d", len(keys))
}
