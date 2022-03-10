package stream

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func newArray(count int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	s := make([]int, count)
	for i := 0; i < count; i++ {
		s[i] = r.Intn(count * 2)
	}
	return s
}

func BenchmarkParallel(b *testing.B) {
	s := newArray(100)

	filter := func(v int) bool {
		time.Sleep(time.Millisecond)
		return v > 0
	}

	b.Run("no parallel", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = NewSlice(s).Parallel(1).Filter(filter)
		}
	})

	b.Run("goroutine 2", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = NewSlice(s).Parallel(2).Filter(filter)
		}

	})

	b.Run("goroutine 4", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = NewSlice(s).Parallel(4).Filter(filter)
		}
	})

	b.Run("goroutine 6", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = NewSlice(s).Parallel(6).Filter(filter)
		}
	})

	b.Run("goroutine 8", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = NewSlice(s).Parallel(8).Filter(filter)
		}
	})

	b.Run("goroutine 10", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = NewSlice(s).Parallel(10).Filter(filter)
		}
	})
}

func TestParallelFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(v int) bool
		want      int
	}{
		{
			name:      "match",
			input:     newArray(100),
			predicate: func(v int) bool { return v < 100 },
		},
		{
			name:      "match",
			input:     newArray(200),
			predicate: func(v int) bool { return v < 100 },
		},
		{
			name:      "match",
			input:     newArray(300),
			predicate: func(v int) bool { return v < 100 },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t,
				NewSliceByOrdered(tt.input).Parallel(10).Filter(tt.predicate).Sort().ToSlice(),
				NewSliceByOrdered(tt.input).Filter(tt.predicate).Sort().ToSlice())
		})
	}
}

func TestParallelMap(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		mapper func(int) int
		want   []int
	}{
		{
			name:   "normal",
			input:  newArray(100),
			mapper: func(i int) int { return i * 2 },
		},
		{
			name:   "empty",
			input:  newArray(100),
			mapper: func(i int) int { return i * 2 },
		},
		{
			name:   "nil",
			input:  newArray(100),
			mapper: func(i int) int { return i * 2 },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t,
				NewSliceByOrdered(tt.input).Parallel(10).Map(tt.mapper).Sort().ToSlice(),
				NewSliceByOrdered(tt.input).Map(tt.mapper).Sort().ToSlice())
		})
	}
}
