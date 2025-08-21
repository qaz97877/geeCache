package geecache

import (
	"geecache/lru"
	"sync"
)

// 注意这个结构体是geecache包内部私有的
type cache struct {
	mutex      sync.Mutex
	lru        *lru.Cache // 指针成员默认初值为nil 所以需要显式初始化
	cacheBytes int64
}

// 方法也都是私有的
// 这个函数并不是单例模式 因为单例模式要求创建的对象是全局唯一的 而这里只保证为每个结构体生成唯一的lru成员
func (c *cache) add(key string, val ByteView) {
	// 注意这里使用读写锁的不可行的 因为当前c中没有lru时会在这个函数进行创建
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, val)
}

func (c *cache) get(key string) (val ByteView, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		// bool类型的默认零值是false
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), true
	}
	return
}
