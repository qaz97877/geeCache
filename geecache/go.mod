module geecache

go 1.24.5

replace geecache/lru => ./lru

replace geecache/consistenthash => ./consistenthash

require google.golang.org/protobuf v1.36.8
