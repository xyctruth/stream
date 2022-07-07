# Stream

[![Build](https://github.com/xyctruth/stream/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/xyctruth/stream/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/xyctruth/stream/branch/main/graph/badge.svg?token=ZHMPMQP0CP)](https://codecov.io/gh/xyctruth/stream)

> [English](./README.md) / [中文](./README-ZH.md)

`Stream` is a stream processing library based on Go 1.18+ Generics. It supports parallel processing of data in the stream.

## Features

-  `Parallel`: Parallel processing of data in the stream, keeping the original order of the elements in the stream
-  `Pipeline`: combine multiple operations to reduce element loops, short-circuiting earlier
-  `Lazy Invocation`: intermediate operations are lazy

## Installation

Requires Go 1.18+ version installed

```go
import "github.com/xyctruth/stream"
```

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

## Type Conversion

Sometimes we need to use `Map` , `Reduce` to convert the type of slice elements, but unfortunately Golang currently does not support structure methods with additional type parameters, all type parameters must be declared in the structure. We work around this with a temporary workaround until Golang supports it.

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

### Parallel Type

- `First`: parallel processing ends as soon as the first return value is obtained. `For: AllMatch, AnyMatch, FindFunc`
- `ALL`: All elements need to be processed in parallel, all return values are obtained, and then the parallel is ended. `For: Map, Filter`
- `Action`: All elements need to be processed in parallel, no return value required. `For: ForEach, Action`

### Parallel Goroutines Number

The number of parallel goroutines has different choices for CPU operations and IO operations. Generally, the number of goroutines does not need to be set larger than the number of CPU cores for CPU operations, while the number of goroutines for IO operations can be set to be much larger than the number of CPU cores.

#### CPU Operations

[BenchmarkParallelByCPU](./benchmark_test.go)

```go
NewSlice(s).Parallel(tt.goroutines).ForEach(func(i int, v int) {
    sort.Ints(newArray(1000)) // Simulate time-consuming CPU operations
})
```
Benchmark with 6 cpu cores
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

#### IO Operations

[BenchmarkParallelByIO](./benchmark_test.go)

```go
NewSlice(s).Parallel(tt.goroutines).ForEach(func(i int, v int) {
    time.Sleep(time.Millisecond) // Simulate time-consuming IO operations
})
```
Benchmark with 6 cpu cores
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

