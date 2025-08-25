package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 哈希函数将bytes映射到uint32
type Hash func(data []byte) uint32

// 包含所有的哈希值
type Map struct {
	hash     Hash
	replicas int            // 虚拟节点的数量
	keys     []int          // 有序的哈希环
	hashMap  map[int]string // key是虚拟节点的哈希值 val是真实节点的名称
}

// 创建一个哈希环结构体
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 可以传入0个或多个真实节点的名称
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// 获取最近的节点方法
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 进行二分查找得到最近的位置
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
