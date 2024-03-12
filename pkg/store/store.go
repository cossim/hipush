package store

type Store interface {
	Init() error
	Get(key string) int64
	Set(key string, value int64)
	Add(key string, value int64)
	Del(key string)
	Close() error
}
