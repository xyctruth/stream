package main

import (
	"fmt"
	"github.com/xyctruth/stream"
	"strconv"
)

func main() {
	// Constraints `any`
	s1 := stream.NewSlice([]string{"d", "a", "b", "c", "a"}).
		Filter(func(s string) bool { return s != "b" }).
		Map(func(s string) string { return "class_" + s }).
		SortFunc(func(s1, s2 string) bool { return s1 < s2 }).
		ToSlice()
	fmt.Println(s1)

	// Constraints `comparable`
	s2 := stream.NewSliceByComparable([]string{"d", "a", "b", "c", "a"}).
		Filter(func(s string) bool { return s != "b" }).
		Map(func(s string) string { return "class_" + s }).
		SortFunc(func(s1, s2 string) bool { return s1 < s2 }).
		Distinct().
		ToSlice()
	fmt.Println(s2)

	//Constraints `constraints.Ordered`
	s3 := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
		Filter(func(s string) bool { return s != "b" }).
		Map(func(s string) string { return "class_" + s }).
		Sort().
		Distinct().
		ToSlice()
	fmt.Println(s3)

	// Parallel
	s4 := stream.NewSliceByOrdered([]string{"d", "a", "b", "c", "a"}).
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
	fmt.Println(s4)

	// Need to convert the type of slice elements.
	s5 := stream.NewSliceByMapping[int, string, string]([]int{1, 2, 3, 4, 5}).
		Filter(func(v int) bool { return v > 3 }).
		Map(func(v int) string { return "mapping_" + strconv.Itoa(v) }).
		Reduce(func(r string, v string) string { return r + v })

	fmt.Println(s5)

}
