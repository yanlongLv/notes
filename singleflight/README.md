singleflight笔记

Group struct
```
type Group struct {
	mu sync.Mutex  //保证并发
	m  map[string]*call // 保存key
}
```
//
```
type call struct {
	wg        sync.WaitGroup //保证只有一个协程执行函数
	val       interface{} // 函数的执行结果
	err       error
	forgotten bool // 是否丢掉
	dups      int //同一个key访问的线程数量
	chans     []chan<- Result
}
```

1.singleflight 是为了解决雪崩缓存穿透开发的工具,不到100行代码解决了一大对问题
//该方法singleflight的具体实现
Do
//该方法singleflight的具体返回chan类型
DoChan
//Forget 丢掉响应的key
Forget
过去应用一些防止缓冲击穿雪崩的方法

1.使用redis的SETNX 在设置过期时间的发放
这种方法的缺陷：
    *没法评估具体的过期时间，应该设置多少合适
    *在线程等待时CPU空转，CPU利用率不高
    *超时控制不好做
