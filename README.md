# Stream [WIP]

## Introduction

Stream is a stream library based on golang 1.18 generics (manipulate slice like java stream)

## Getting Started

### Constraints `any`

```go
s1 := stream.NewSlice([]string{"d", "a", "b", "c", "a"}).
    Filter(func(s string) bool { return s != "b" }).
    Map(func(s string) string { return "class_" + s }).
    SortFunc(func(s1, s2 string) bool { return s1 < s2 }).
    Slice()

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
    Slice()

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
    Slice()
    fmt.Println(s3)

fmt.Println(s3)

//output: [class_a class_c class_d]
```

