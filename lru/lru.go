package lru

import "container/list"

type Cache struct {
	maxBytes int64
	nbytes   int64
	list     *list.List
	cache    map[string]*list.Element
	callBack func(key string, val Value)
}

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
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(val.Len()) - int64(kv.val.Len())
		kv.val = val
	} else {
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
