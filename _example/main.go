package main

import (
	"fmt"
	"github.com/xyctruth/stream"
)

func main() {
	s1 := stream.NewSlice([]string{"d", "a", "b", "c", "a"}).
		Filter(func(s string) bool { return s != "b" }).
		Map(func(s string) string { return "class_" + s }).
		SortFunc(func(s1, s2 string) bool { return s1 < s2 }).
		ToSlice()
	fmt.Println(s1)

	s2 := stream.NewSliceByComparable([]string{"d", "a", "b", "c", "a"}).
		Filter(func(s string) bool { return s != "b" }).
		Map(func(s string) string { return "class_" + s }).
		SortFunc(func(s1, s2 string) bool { return s1 < s2 }).
		Distinct().
		ToSlice()
	fmt.Println(s2)

	s3 := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
		Filter(func(s string) bool { return s != "b" }).
		Map(func(s string) string { return "class_" + s }).
		Sort().
		Distinct().
		ToSlice()
	fmt.Println(s3)

}
