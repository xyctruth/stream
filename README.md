# Stream [WIP]

[![Build](https://github.com/xyctruth/stream/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/xyctruth/stream/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/xyctruth/stream/branch/main/graph/badge.svg?token=ZHMPMQP0CP)](https://codecov.io/gh/xyctruth/stream)

## Introduction

Stream is a stream library based on golang 1.18 generics, support parallel. (manipulate slice like java stream)

## Getting Started

### Basic

```go
s := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    Sort().
    Distinct().
    ToSlice()
```

### Constraints

`any`constraint means that any type of slice can be accepted, can't use `==` `!=` `>` `<` to compare elements, so you can't use `Sort()`, `Find()`... , but you can use `SortFunc(fn)`, `FindFunc(fn)`... instead of

```go
stream.NewSlice([]int{1, 2, 3, 7, 1})
```

`comparable` constraint accepts types that can use `==` `!=`, but still can't use `>` `<` to compare elements, so you can't use `SortFunc(fn)`... , but you can use `SortFunc (fn)`... instead of

```go
stream.NewSliceByComparable([]int{1, 2, 3, 7, 1})
```

`constraints.Ordered`constraint accepts types that can use `==` `!=` `>` `<`, so can use all functions

```go
stream.NewSliceByOrdered([]int{1, 2, 3, 7, 1})
```

### Parallel

The `Parallel` function accepts a `goroutines int` parameter. If goroutines>1, create `goroutines` goroutines to execute, otherwise close Parallel, the stream closes Parallel by default.

```go
s := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Parallel(10).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    Sort().
    Distinct().
    ToSlice()
```
