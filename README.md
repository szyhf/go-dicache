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