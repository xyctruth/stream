# Stream: 基于 Go 1.18+ 泛型的流式处理库 (支持并行流)

[![Build](https://github.com/xyctruth/stream/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/xyctruth/stream/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/xyctruth/stream/branch/main/graph/badge.svg?token=ZHMPMQP0CP)](https://codecov.io/gh/xyctruth/stream)

> [English](./README.md) / [中文](./README-ZH.md)

Stream 是一个基于 Go 1.18+ 泛型的流式处理库, 它支持并行处理流中的数据. 并行流会将元素平均划分多个的分区, 并创建相同数量的 goroutine 执行, 并且会保证处理完成后流中元素保持原始顺序.

## 入门

```go
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

### 并行 goroutines

开启并行 goroutine 数量在面对 CPU 操作与 IO 操作有着不同的选择。 一般面对 CPU 操作时 goroutine 数量不需要设置大于 CPU 核心数，而 IO 操作时 goroutine 数量可以设置远远大于 CPU 核心数.

#### CPU 操作

[BenchmarkParallelByCPU](./parallel_test.go)

```go
NewSlice(s).Parallel(goroutines).ForEach(func(i int, v int) {
    sort.Ints(newArray(1000)) // 模拟 CPU 耗时操作
})
```

```go
go test -run=^$ -benchtime=5s -cpu=6  -bench=^BenchmarkParallelByCPU

goarch: amd64
pkg: github.com/xyctruth/stream
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
BenchmarkParallelByCPU/no_parallel(0)-6         	     717	   9183119 ns/op
BenchmarkParallelByCPU/goroutines(2)-6          	    1396	   4303113 ns/op
BenchmarkParallelByCPU/goroutines(4)-6          	    2539	   2388197 ns/op
BenchmarkParallelByCPU/goroutines(6)-6          	    2932	   2159407 ns/op
BenchmarkParallelByCPU/goroutines(8)-6          	    2334	   2577405 ns/op
BenchmarkParallelByCPU/goroutines(10)-6         	    2649	   2352926 ns/op
```

#### IO 操作

[BenchmarkParallelByIO](./parallel_test.go)

```go
NewSlice(s).Parallel(goroutines).ForEach(func(i int, v int) {
    time.Sleep(time.Millisecond) // 模拟 IO 耗时操作
})
```

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
