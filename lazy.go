package diCache

import (
	"sync"
)

// 缓存的行级元素
type lazyCacheElem struct {
	// 缓存值
	val interface{}
	// 重建结果异常
	err error
	// 等待标志
	ready chan struct{}
}

// 具备延迟加载特性的kv缓存
// 使用时要设置延迟加载方法。
type LazyCache struct {
	// 全缓存锁
	mx *sync.Mutex

	// 全缓存
	memo map[interface{}]*lazyCacheElem

	// 重构方法
	newElemFunc func(key interface{}) (interface{}, error)
}

// 获取指定key的缓存，如果缓存不存在则调用重构方法尝试重构。
func (this *LazyCache) Get(key interface{}) (interface{}, error) {
	this.mx.Lock()
	elem, ok := this.memo[key]
	if !ok {
		// 第一次访问，重建
		elem = &lazyCacheElem{ready: make(chan struct{})}
		this.memo[key] = elem
		this.mx.Unlock()

		// 加载数据
		elem.val, elem.err = this.newElemFunc(key)

		close(elem.ready)
	} else {
		// 已有缓存
		this.mx.Unlock()
		// 等待加载
		<-elem.ready
	}
	return elem.val, elem.err
}

// 删除指定key的缓存
func (this *LazyCache) Del(key interface{}) {
	this.mx.Lock()
	delete(this.memo, key)
	this.mx.Unlock()
}

// 写入
func (this *LazyCache) Put(key, val interface{}) {
	this.mx.Lock()
	elem := &lazyCacheElem{ready: make(chan struct{})}
	this.memo[key] = elem
	this.mx.Unlock()
	elem.val = val
	close(elem.ready)
}
