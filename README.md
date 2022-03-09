# Stream [WIP]

[![Build](https://github.com/xyctruth/stream/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/xyctruth/stream/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/xyctruth/stream/branch/main/graph/badge.svg?token=ZHMPMQP0CP)](https://codecov.io/gh/xyctruth/stream)

## Introduction

Stream is a stream library based on golang 1.18 generics (manipulate slice like java stream)

## Getting Started

### Constraints `any`

```go
s1 := stream.NewSlice([]string{"d", "a", "b", "c", "a"}).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    SortFunc(func(s1, s2 string) bool { return s1 < s2 }).
    ToSlice()

fmt.Println(s1)

//output: [class_a class_a class_c class_d]
```

### Constraints `comparable`

```go
s2 := stream.NewSliceByComparable([]string{"d", "a", "b", "c", "a"}).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    SortFunc(func(s1, s2 string) bool { return s1 < s2 }).
    Distinct().
    ToSlice()

fmt.Println(s2)

//output:[class_a class_c class_d]
```

### Constraints `constraints.Ordered`

```go
s3 := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    Sort().
    Distinct().
    ToSlice()

fmt.Println(s3)

//output: [class_a class_c class_d]
```

