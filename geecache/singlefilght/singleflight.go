package singleflight

import "sync"

// 在进行中或已结束的请求
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// 管理不同key的请求
type Group struct {
	mu sync.Mutex
	m  map[string]*call // 请求缓存 保证同一时间同一个key只被查询一次
}

// 针对同一个key 无论Do被调用多少次fn只被调用一次
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() // 若请求正在进行 则等待
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)  // 发起请求前加锁
	g.m[key] = c // 添加到哈希表 表明key已在处理
	g.mu.Unlock()

	c.val, c.err = fn() // 调用fn发起请求
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	delete(g.m, key) // 更新哈希表
	g.mu.Unlock()

	return c.val, c.err
}
