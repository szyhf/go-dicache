# 介绍

简单的数据缓存结构，功能比redis弱很多，仅仅是简单，需求丰富的话还是用redis比较靠谱。

```shell
go get github.com/szyhf/go-dicache
```

## LazyCache

具备延迟加载特性的kv缓存，使用时要设置延迟加载方法。

``` go
import (
    "fmt"
    dicache "github.com/szyhf/go-dicache"
)
lazyCache := dicache.NewLazyCache(func(key interface{}) (interface{}, error){
    // 自己实现延迟加载的重构函数，会在第一次访问这个key的时候执行。
    if key == "Hello"{
        return "World",nil
    }else{
        return nil,fmt.Errorf("Foo")
    }
})
world,err:=lazyCache.Get("Hello")
fmt.Println(world,err)
// world = "World", err = nil
bar,err:=lazyCache.Get("Foo")
fmt.Println(bar,err)
// bar = nil, err = Error("Foo")
```

组合式用法（在没有泛型的情况下尽可能让人用起来苏胡一点）:

```go
// 可以是自己定义的任何可以作为map key的类型
type MyKeyType interface{}
// 可是是自己定义的任何可以作为map val的类型
type MyValType interface{}

type MyLazyCache struct{
	*dicache.LazyCache
}

func NewMyLazyCache() *MyLazyCache{
	lazyCache := dicache.NewLazyCache(func(key interface{}) (interface{}, error){
    	// 自己实现延迟加载的重构函数，会在第一次访问这个key的时候执行。
    	if key == "Hello"{
        	return "World",nil
    	}else{
        	return nil,fmt.Errorf("Foo")
    	}
	})
	return &MyLazyCache{
		LazyCache:lazyCache,
	}
}

// 重载Get方法，让它用起来符合预期（模拟泛型）
func (this*MyLazyCache)Get(key MyKeyType)(*MyValType,error){
	res,err:=this.LazyCache.Get(key)
	return res.(*MyValType),err
}

```

> 一般来说重载Get是最重要的，其他的自己看需求酌情处理就好。