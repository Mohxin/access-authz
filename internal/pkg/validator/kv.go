package validator

import "sync"

// KV struct using sync.Map with generics
type KV[K comparable, V any] struct {
	hash sync.Map
}

func NewKV[K comparable, V any]() *KV[K, V] {
	return &KV[K, V]{hash: sync.Map{}}
}

// Set a value in the cache
func (kv *KV[K, V]) Set(key K, value V) {
	kv.hash.Store(key, value)
}

// Get a value by key from the cache
func (kv *KV[K, V]) Get(key K) (V, bool) {
	val, ok := kv.hash.Load(key)
	if ok {
		v, ok := val.(V)
		if ok {
			return v, ok
		}
	}
	var zero V // Return zero value if key not found
	return zero, false
}

// Contains checks if a key exists in the cache
func (kv *KV[K, V]) Contains(key K) bool {
	_, exists := kv.Get(key)
	return exists
}

// Delete a value by key from the cache
func (kv *KV[K, V]) Delete(key K) {
	kv.hash.Delete(key)
}

// List all values in the cache
func (kv *KV[K, V]) List() map[K]V {
	m := make(map[K]V)
	kv.hash.Range(func(key, value interface{}) bool {
		k, ok := key.(K)
		if !ok {
			return false
		}

		val, ok := value.(V)
		if !ok {
			return false
		}

		m[k] = val
		return true
	})

	return m
}

func (kv *KV[K, V]) Filter(fn func(K, V) bool) {
	kv.hash.Range(func(key, value interface{}) bool {
		k, ok := key.(K)
		if !ok {
			return false
		}

		val, ok := value.(V)
		if !ok {
			return false
		}

		if !fn(k, val) {
			return true
		}

		return true
	})
}

func (kv *KV[K, V]) Values() []V {
	var values []V
	kv.hash.Range(func(_, value interface{}) bool {
		v, ok := value.(V)
		if !ok {
			return false
		}
		values = append(values, v)
		return true
	})
	return values
}

func (kv *KV[K, V]) Keys() []K {
	var keys []K
	kv.hash.Range(func(key, _ interface{}) bool {
		k, ok := key.(K)
		if !ok {
			return false
		}
		keys = append(keys, k)
		return true
	})
	return keys
}
