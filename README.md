# Stream

[![Build](https://github.com/xyctruth/stream/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/xyctruth/stream/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/xyctruth/stream/branch/main/graph/badge.svg?token=ZHMPMQP0CP)](https://codecov.io/gh/xyctruth/stream)

> [English](./README.md) / [中文](./README-ZH.md)

## Introduction

Stream is a stream library based on golang 1.18+ generics, support parallel. (manipulate slice like java stream)

## Getting Started

```go
s := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    Sort().
    Distinct().
    ToSlice()
```

### Constraints

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

### Parallel

The `Parallel` function accept a `goroutines int` parameter. If goroutines>1, open Parallel , otherwise close Parallel, the stream Parallel is off by default.

Parallel will divide the elements in the stream into multiple partitions equally, and create the same number of goroutine to execute, and the original order of the elements of the stream will be guaranteed.

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
