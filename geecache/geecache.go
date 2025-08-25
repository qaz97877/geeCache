package geecache

import (
	"fmt"
	"log"
	"sync"
)

// 用于加载数据的回调函数
// geeCache不应该实现具体的加载方法 而是应该让用户传入加载函数
type Getter interface {
	Get(key string) ([]byte, error)
}

// 函数适配器
// 因为普通的函数是不具有Get方法的 所以无法通过Getter(simpleFunc)这种形式强转
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 一个集群节点的格式
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
}

// 全局变量
var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 情况1：在当前的缓存中找到了key 直接返回data
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (val ByteView, err error) {
	// 情况2：当前缓存中没有找到key 与远程节点交互 返回缓存值
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if bytes, err := g.getFromPeer(peer, key); err == nil {
				return bytes, nil
			}
			log.Println("[GeeCache] Failed to get from peer", err)
		}
	}
	// 情况3：在当前缓存中没找到key 且在远程节点也未找到
	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	val := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, val)
	return val, nil
}

func (g *Group) populateCache(key string, val ByteView) {
	g.mainCache.add(key, val)
}

// 注册PeerPicker用于选择远程peer
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}
