package diCache

import "sync"

type Cacher interface {
	Get(key interface{}) (interface{}, error)
	Del(key interface{})
	Put(key, val interface{})
	IsExist(key interface{}) bool
}

// 新建一个LazyCache
func NewLazyCache(newElemFunc func(key interface{}) (interface{}, error)) *LazyCache {
	c := &LazyCache{
		memo:        &sync.Map{},
		newElemFunc: newElemFunc,
	}
	return c
}
