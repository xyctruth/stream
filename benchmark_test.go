package stream

import (
	"fmt"
	"testing"
)

func BenchmarkShortCircuiting(b *testing.B) {
	tests := []struct {
		name  string
		count int
	}{
		{count: 100},
		{count: 200},
		{count: 300},
		{count: 400},
		{count: 500},
		{count: 1000},
		{count: 2000},
	}
	for _, tt := range tests {
		b.Run(fmt.Sprintf("%s(%d)", "count:", tt.count), func(b *testing.B) {
			s := newArray(tt.count)
			s[0] = 101
			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				_ = NewSlice(s).
					Filter(func(v int) bool { return true }).
					Map(func(v int) int {
						return v
					}).
					AllMatch(func(v int) bool { return v < 100 })
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
