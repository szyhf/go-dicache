package diCache

import "sync"

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
	// 全缓存
	memo *sync.Map

	// 重构方法
	newElemFunc func(key interface{}) (interface{}, error)
}

// 获取指定key的缓存，如果缓存不存在则调用重构方法尝试重构。
func (this *LazyCache) Get(key interface{}) (interface{}, error) {
	e, ok := this.memo.Load(key)
	elem := e.(*lazyCacheElem)
	if !ok {
		// 第一次访问，重建
		elem = &lazyCacheElem{ready: make(chan struct{})}
		this.memo.Store(key, elem)

		// 加载数据
		elem.val, elem.err = this.newElemFunc(key)

		close(elem.ready)
	} else {
		// 等待加载
		<-elem.ready
	}
	return elem.val, elem.err
}

// 删除指定key的缓存
func (this *LazyCache) Del(key interface{}) {
	this.memo.Delete(key)
}

// 写入
func (this *LazyCache) Put(key, val interface{}) {
	elem := &lazyCacheElem{ready: make(chan struct{})}
	elem.val = val
	this.memo.Store(key, elem)
	close(elem.ready)
}

// 存在性判断
func (this *LazyCache) IsExist(key interface{}) bool {
	_, ok := this.memo.Load(key)
	return ok
}
