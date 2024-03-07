package store

import (
	"strconv"
	"sync"
	"sync/atomic"
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: sync.Map{},
	}
}

type MemoryStore struct {
	data sync.Map
}

func (m *MemoryStore) Init() error {
	return nil
}

func (m *MemoryStore) Get(key string) int64 {
	if val, ok := m.data.Load(key); ok {
		return val.(int64)
	}
	return 0
}

func (m *MemoryStore) Set(key string, value int64) {
	m.data.Store(key, value)
}

func (m *MemoryStore) Add(key string, value int64) {
	for {
		oldValue, loaded := m.data.LoadOrStore(key, int64(0))
		if loaded {
			// 将加载的值转换为int64类型
			old := oldValue.(int64)
			// 使用原子操作更新值
			newValue := atomic.AddInt64(&old, value)
			// 如果原子操作成功，则返回
			if m.data.CompareAndSwap(key, oldValue, newValue) {
				return
			}
		} else {
			return
		}
	}
}

func (m *MemoryStore) Del(key string) {
	m.data.Delete(key)
}

func (m *MemoryStore) Close() error {
	// MemoryStore doesn't need to be closed, so just return nil
	return nil
}

func (m *MemoryStore) String() string {
	var str string
	m.data.Range(func(key, value interface{}) bool {
		str += key.(string) + ": " + strconv.FormatInt(value.(int64), 10) + "\n"
		return true
	})
	return str
}
