package geecache

// 用于封装缓存值的只读数据结构
// 实现只读的方法是用值传递类型与方法绑定
// 支持任意类型的数据
type ByteView struct {
	b []byte
}

// 得到缓存的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// 得到缓存的副本
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// 将缓存转换为字符串
func (v ByteView) String() string {
	return string(v.b)
}

// 生成副本的实现 深拷贝
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
