# Stream [WIP]

[![Build](https://github.com/xyctruth/stream/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/xyctruth/stream/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/xyctruth/stream/branch/main/graph/badge.svg?token=ZHMPMQP0CP)](https://codecov.io/gh/xyctruth/stream)

## Introduction

Stream is a stream library based on golang 1.18 generics, support parallel. (manipulate slice like java stream)

## Getting Started

### Constraints 

#### any

```go
s1 := stream.NewSlice([]string{"d", "a", "b", "c", "a"}).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    SortFunc(func(s1, s2 string) bool { return s1 < s2 }).
    ToSlice()
```

#### comparable

```go
s2 := stream.NewSliceByComparable([]string{"d", "a", "b", "c", "a"}).
    Distinct().
    ToSlice()
```

#### constraints.Ordered

```go
s3 := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Sort().
    Distinct().
    ToSlice()
```

### Parallel

```go
s4 := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Parallel(10).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    Sort().
    Distinct().
    ToSlice()
```
