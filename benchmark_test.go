package stream

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func BenchmarkParallelByCPU(b *testing.B) {
	tests := []struct {
		name       string
		goroutines int
		action     func(int, int)
	}{
		{name: "no Parallel", goroutines: 0},
		{name: "Goroutines", goroutines: 2},
		{name: "Goroutines", goroutines: 4},
		{name: "Goroutines", goroutines: 6},
		{name: "Goroutines", goroutines: 8},
		{name: "Goroutines", goroutines: 10},
	}
	s := newArray(100)

	for _, tt := range tests {
		b.Run(fmt.Sprintf("%s(%d)", tt.name, tt.goroutines), func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				NewSlice(s).Parallel(tt.goroutines).ForEach(func(i int, v int) {
					sort.Ints(newArray(1000)) // Simulate time-consuming CPU operations
				})
			}
		})
	}
}

func BenchmarkParallelByIO(b *testing.B) {
	tests := []struct {
		name       string
		goroutines int
		action     func(int, int)
	}{
		{name: "no Parallel", goroutines: 0},
		{name: "Goroutines", goroutines: 2},
		{name: "Goroutines", goroutines: 4},
		{name: "Goroutines", goroutines: 6},
		{name: "Goroutines", goroutines: 8},
		{name: "Goroutines", goroutines: 10},
		{name: "Goroutines", goroutines: 50},
		{name: "Goroutines", goroutines: 100},
	}
	s := newArray(100)

	for _, tt := range tests {
		b.Run(fmt.Sprintf("%s(%d)", tt.name, tt.goroutines), func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				NewSlice(s).Parallel(tt.goroutines).ForEach(func(i int, v int) {
					time.Sleep(time.Millisecond) // Simulate time-consuming IO operations
				})
			}
		})
	}
}

var funcVariable = mapperFunc

func mapperFunc(v int) int { return v * 2 }

// TODO optimize func variable 3.287 ns/op, other 1.553 ns/op
func BenchmarkPipelineFuncVariable(b *testing.B) {
	tests := []struct {
		name       string
		goroutines int
		stage      func(index int, e int) (isReturn bool, isComplete bool, ret int)
	}{
		{
			name: "func variable",
			stage: func(index int, e int) (isReturn bool, isComplete bool, ret int) {
				return true, false, funcVariable(e)
			},
		},
		{
			name: "normal func",
			stage: func(index int, e int) (isReturn bool, isComplete bool, ret int) {
				return true, false, mapperFunc(e)
			},
		},
		{
			name: "anonymous func",
			stage: func(index int, e int) (isReturn bool, isComplete bool, ret int) {
				return true, false, func(v int) int { return v * 2 }(e)
			},
		},
	}
	for _, tt := range tests {
		b.Run(fmt.Sprintf("%s(%d)", tt.name, tt.goroutines), func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				tt.stage(0, 0)
			}
		})
	}
}
