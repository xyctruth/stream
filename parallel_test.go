package stream

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
				_ = NewSlice(s).Parallel(tt.goroutines).ForEach(func(i int, v int) {
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
				_ = NewSlice(s).Parallel(tt.goroutines).ForEach(func(i int, v int) {
					time.Sleep(time.Millisecond) // Simulate time-consuming IO operations
				})
			}
		})
	}
}

func BenchmarkParallelAForEach(b *testing.B) {
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
	s := newArray(10000000)

	for _, tt := range tests {
		b.Run(fmt.Sprintf("%s(%d)", tt.name, tt.goroutines), func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = NewSlice(s).Parallel(tt.goroutines).ForEach(func(i int, v int) {
				})
			}
		})
	}
}

func BenchmarkParallelAMap(b *testing.B) {
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
	s := newArray(10000000)

	for _, tt := range tests {
		b.Run(fmt.Sprintf("%s(%d)", tt.name, tt.goroutines), func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = NewSlice(s).Parallel(tt.goroutines).Map(func(v int) int {
					return v * 2
				})
			}
		})
	}
}

func TestPartition(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		goroutine int
		want1     []part
		want2     int
	}{
		{
			name:      "case",
			input:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			goroutine: 2,
			want1: []part{
				{
					low:  0,
					high: 6,
				},
				{
					low:  6,
					high: 13,
				},
			},
		},
		{
			name:      "case",
			input:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			goroutine: 3,
			want1: []part{
				{
					low:  0,
					high: 4,
				},
				{
					low:  4,
					high: 8,
				},
				{
					low:  8,
					high: 13,
				},
			},
		},
		{
			name:      "case",
			input:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			goroutine: 5,
			want1: []part{
				{
					low:  0,
					high: 2,
				},
				{
					low:  2,
					high: 4,
				},
				{
					low:  4,
					high: 7,
				},
				{
					low:  7,
					high: 10,
				},
				{
					low:  10,
					high: 13,
				},
			},
		},
		{
			name:      "case",
			input:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			goroutine: 10,
			want1: []part{
				{
					low:  0,
					high: 1,
				},
				{
					low:  1,
					high: 2,
				},
				{
					low:  2,
					high: 3,
				},
				{
					low:  3,
					high: 4,
				},
				{
					low:  4,
					high: 5,
				},
				{
					low:  5,
					high: 6,
				},
				{
					low:  6,
					high: 7,
				},
				{
					low:  7,
					high: 9,
				},
				{
					low:  9,
					high: 11,
				},
				{
					low:  11,
					high: 13,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := partition(tt.input, tt.goroutine)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
