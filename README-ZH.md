# Stream [WIP]

[![Build](https://github.com/xyctruth/stream/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/xyctruth/stream/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/xyctruth/stream/branch/main/graph/badge.svg?token=ZHMPMQP0CP)](https://codecov.io/gh/xyctruth/stream)

> [English](./README.md) / [中文](./README-ZH.md)

## 简介

Stream 是一个基于 golang 1.18+ 泛型的流库，支持并行。(像Java流一样操作切片)

## 入门

### 基础

```go
s := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    Sort().
    Distinct().
    ToSlice()
```

### 约束

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

### 并行

`Parallel` 函数接收一个 `goroutines int` 参数. 如果 goroutines>1 则创建相同数量的goroutines执行, 否则关闭并行, 默认流是关闭并行的。

```go
s := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
    Parallel(10).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    Sort().
    Distinct().
    ToSlice()
```
