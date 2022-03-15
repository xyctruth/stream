# Stream: A Stream library based on Go 1.18+ Generics (Support Parallel Stream)

[![Build](https://github.com/xyctruth/stream/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/xyctruth/stream/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/xyctruth/stream/branch/main/graph/badge.svg?token=ZHMPMQP0CP)](https://codecov.io/gh/xyctruth/stream)

> [English](./README.md) / [中文](./README-ZH.md)

Stream is a Stream library based on Go 1.18+ Generics. It supports parallel processing of data in the stream. The parallel stream will divide the elements into multiple partitions equally, and create the same number of goroutine for execute, and will ensure that the elements in the stream remain in the original order after the processing is complete.

## Getting Started

```go
s := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    Sort().
    Distinct().
    ToSlice()
```

## Type Constraints

`any` accepts elements of any type, so you cannot use `==` `!=` `>` `<` to compare elements, which will prevent you from using Sort(), Find()... functions, but you can use SortFunc(fn), FindFunc(fn)... instead

```go
stream.NewSlice([]int{1, 2, 3, 7, 1})
```

`comparable` accepts type can use `==` `!=` to compare elements, but still can't use `>` `<` to compare elements, so you can't use Sort(), Min()... functions, but you can use SortFunc(fn), MinFunc()... instead

```go
stream.NewSliceByComparable([]int{1, 2, 3, 7, 1})
```

`constraints.Ordered` accepts types that can use `==` `!=` `>` `<`  to compare elements, so can use all functions

```go
stream.NewSliceByOrdered([]int{1, 2, 3, 7, 1})
```

## Parallel

The `Parallel` function accept a `goroutines int` parameter. If goroutines>1, open Parallel , otherwise close Parallel, the stream Parallel is off by default.

Parallel will divide the elements in the stream into multiple partitions equally, and create the same number of goroutine to execute, and it will ensure that the elements in the stream remain in the original order after processing is complete.

```go
s := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Parallel(10).
    Filter(func(s string) bool {
    // some time-consuming operations
    return s != "b"
    }).
    Map(func(s string) string {
    // some time-consuming operations
    return "class_" + s
    }).
    ForEach(
    func(index int, s string) {
    // some time-consuming operations
    },
    ).ToSlice()
```

### Parallel goroutines

The number of parallel goroutines has different choices for CPU operations and IO operations. Generally, the number of goroutines does not need to be set larger than the number of CPU cores for CPU operations, while the number of goroutines for IO operations can be set to be much larger than the number of CPU cores.

#### CPU Operations

[BenchmarkParallelByCPU](./parallel_test.go)

```go
NewSlice(s).Parallel(tt.goroutines).ForEach(func(i int, v int) {
    sort.Ints(newArray(1000)) // Simulate time-consuming CPU operations
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

#### IO Operations

[BenchmarkParallelByIO](./parallel_test.go)

```go
NewSlice(s).Parallel(tt.goroutines).ForEach(func(i int, v int) {
    time.Sleep(time.Millisecond) // Simulate time-consuming IO operations
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

