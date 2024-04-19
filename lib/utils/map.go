package utils

import "sync"

// 并发安全map
type MapSync[T any] struct {
	m    map[string]T
	lock sync.RWMutex
}

func (this *MapSync[T]) Set(key string, val T) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.m[key] = val
}

func (this *MapSync[T]) Get(key string) T {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.m[key]
}
func (this *MapSync[T]) Keys() []string {
	this.lock.RLock()
	defer this.lock.RUnlock()
	keys := make([]string, 0, len(this.m))
	for k := range this.m {
		keys = append(keys, k)
	}
	return keys
}
func (this *MapSync[T]) Values() []T {
	this.lock.RLock()
	defer this.lock.RUnlock()
	values := make([]T, 0, len(this.m))
	for k := range this.m {
		values = append(values, this.m[k])
	}
	return values
}
func NewMap[T any]() *MapSync[T] {
	return &MapSync[T]{m: make(map[string]T)}
}
