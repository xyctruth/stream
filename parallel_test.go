package stream

import (
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
		//time.Sleep(time.Millisecond)
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
