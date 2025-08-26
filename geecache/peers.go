package geecache

import pb "geecache/geecachepb"

// 这个抽象接口用于选择一个peer
// 抽象接口用于实现类似只有成员函数的虚基类而且不是通过显式继承实现
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 这个抽象接口用于实现向相应的peer发送key
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
