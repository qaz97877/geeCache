package lru

import "container/list"

type Cache struct {
	maxBytes int64      // 一个Cache中可存放的最大字节数
	nbytes   int64      // 当前以存放的字节数 包含key和value
	list     *list.List // 一个双链表用于淘汰长期未使用的key
	cache    map[string]*list.Element
	callBack func(key string, val Value) // 移除key时的回调函数 这里只是用于说明回调的用法
}

// 封装键值对
type entry struct {
	key string
	val Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, cb func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		list:     list.New(),
		cache:    make(map[string]*list.Element),
		callBack: cb,
	}
}

func (c *Cache) Get(key string) (val Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.val, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.list.Back()
	if ele != nil {
		c.list.Remove(ele)
		// 类型断言
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.val.Len())
		if c.callBack != nil {
			c.callBack(kv.key, kv.val)
		}
	}
}

func (c *Cache) Add(key string, val Value) {
	if ele, ok := c.cache[key]; ok {
		// 如果已经在缓存中 则更新缓存 调整链表
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(val.Len()) - int64(kv.val.Len())
		kv.val = val
	} else {
		// 如果不在缓存中 则新建缓存项
		ele := c.list.PushFront(&entry{key, val})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(val.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.list.Len()
}
