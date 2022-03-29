# Stream

[![Build](https://github.com/xyctruth/stream/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/xyctruth/stream/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/xyctruth/stream/branch/main/graph/badge.svg?token=ZHMPMQP0CP)](https://codecov.io/gh/xyctruth/stream)

> [English](./README.md) / [中文](./README-ZH.md)

Stream 是一个基于 Go 1.18+ 泛型的流式处理库, 它支持并行处理流中的数据.

## 特性

- [x] 并行流: 保持并行处理完成后流中元素保持原始顺序
- [x] 流水线: 组合多个操作以减少元素循环，更早地短路
- [x] 惰性调用: 中间操作是惰性的

## 安装

需要安装 Go 1.18+ 版本

```go
import "github.com/xyctruth/stream"
```

## 入门

```go
// constraints.Ordered
s := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    Sort().
    Distinct().
    ToSlice()
```

## 类型约束

`any` 接受任何类型的元素, 所以不能使用 `==` `!=` `>` `<` 比较元素, 导致你不能使用 Sort(), Find()...等函数 ,但是你可以使用 SortFunc(fn), FindFunc(fn)... 代替

```go
stream.NewSlice([]int{1, 2, 3, 7, 1})
```

`comparable` 接收的类型可以使用 `==` `!=` 比较元素, 但仍然不能使用 `>` `<` 比较元素, 因此你不能使用 Sort(), Min()...等函数 ,但是你可以使用 SortFunc(fn), MinFunc()... 代替

```go
stream.NewSliceByComparable([]int{1, 2, 3, 7, 1})
```

`constraints.Ordered` 接收的类型可以使用 `==` `!=` `>` `<`, 所以可以使用所有的函数

```go
stream.NewSliceByOrdered([]int{1, 2, 3, 7, 1})
```

## 类型转换

有些时候我们需要使用 `Map` ,`Reduce` 转换切片元素的类型,但是很遗憾目前 Golang 并不支持结构体的方法有额外的类型参数,所有类型参数必须在结构体中声明。在 Golang 支持之前我们暂时使用临时方案解决这个问题。

```go
// SliceMappingStream  Need to convert the type of slice elements.
// - E elements type
// - MapE map elements type
// - ReduceE reduce elements type
type SliceMappingStream[E any, MapE any, ReduceE any] struct {
    SliceStream[E]
}

s := stream.NewSliceByMapping[int, string, string]([]int{1, 2, 3, 4, 5}).
    Filter(func(v int) bool { return v >3 }).
    Map(func(v int) string { return "mapping_" + strconv.Itoa(v) }).
    Reduce(func(r string, v string) string { return r + v })
```

## 并行

`Parallel` 函数接收一个 `goroutines int` 参数. 如果 goroutines>1 则开启并行, 否则关闭并行, 默认流是关闭并行的。

并行会将流中的元素平均划分多个的分区, 并创建相同数量的 goroutine 执行, 并且会保证处理完成后流中元素保持原始顺序.

```go
s := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Parallel(10).
    Filter(func(s string) bool {
    // 一些耗时操作
    return s != "b"
    }).
    Map(func(s string) string {
    // 一些耗时操作
    return "class_" + s
    }).
    ForEach(
    func(index int, s string) {
    // 一些耗时操作
    },
    ).ToSlice()
```

### 并行类型

- `First`: 一旦获得第一个返回值，并行处理就结束. `For: AllMatch, AnyMatch, FindFunc`
- `ALL`: 所有元素都需要并行处理，得到所有返回值，然后并行结束. `For: Map, Filter`
- `Action`: 所有元素需要并行处理，不需要返回值. `For: ForEach, Action`

### 并行 Goroutines 数量

开启并行 goroutine 数量在面对 CPU 操作与 IO 操作有着不同的选择。 一般面对 CPU 操作时 goroutine 数量不需要设置大于 CPU 核心数，而 IO 操作时 goroutine 数量可以设置远远大于 CPU 核心数.

#### CPU 操作

[BenchmarkParallelByCPU](./parallel_test.go)

```go
NewSlice(s).Parallel(goroutines).ForEach(func(i int, v int) {
    sort.Ints(newArray(1000)) // 模拟 CPU 耗时操作
})
```
使用6个cpu核心进行基准测试
```go
go test -run=^$ -benchtime=5s -cpu=6  -bench=^BenchmarkParallelByCPU

goos: darwin
goarch: amd64
pkg: github.com/xyctruth/stream
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
BenchmarkParallelByCPU/no_Parallel(0)-6         	     710	   8265106 ns/op
BenchmarkParallelByCPU/goroutines(2)-6          	    1387	   4333929 ns/op
BenchmarkParallelByCPU/goroutines(4)-6          	    2540	   2361783 ns/op
BenchmarkParallelByCPU/goroutines(6)-6          	    3024	   2100158 ns/op
BenchmarkParallelByCPU/goroutines(8)-6          	    2347	   2531435 ns/op
BenchmarkParallelByCPU/goroutines(10)-6         	    2622	   2306752 ns/op
```

#### IO 操作

[BenchmarkParallelByIO](./parallel_test.go)

```go
NewSlice(s).Parallel(goroutines).ForEach(func(i int, v int) {
    time.Sleep(time.Millisecond) // 模拟 IO 耗时操作
})
```
使用6个cpu核心进行基准测试
```go
go test -run=^$ -benchtime=5s -cpu=6  -bench=^BenchmarkParallelByIO

goos: darwin
goarch: amd64
pkg: github.com/xyctruth/stream
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
BenchmarkParallelByIO/no_parallel(0)-6          	      52	 102023558 ns/op
BenchmarkParallelByIO/goroutines(2)-6           	     100	  55807303 ns/op
BenchmarkParallelByIO/goroutines(4)-6           	     214	  27868725 ns/op
BenchmarkParallelByIO/goroutines(6)-6           	     315	  18925789 ns/op
BenchmarkParallelByIO/goroutines(8)-6           	     411	  14439700 ns/op
BenchmarkParallelByIO/goroutines(10)-6          	     537	  11164758 ns/op
BenchmarkParallelByIO/goroutines(50)-6          	    2629	   2310602 ns/op
BenchmarkParallelByIO/goroutines(100)-6         	    5094	   1221887 ns/op
```
