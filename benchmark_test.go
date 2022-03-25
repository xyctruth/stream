package stream

import (
	"fmt"
	"testing"
)

func BenchmarkPipeline(b *testing.B) {
	tests := []struct {
		name       string
		goroutines int
		action     func(int, int)
	}{
		{name: "no Parallel", goroutines: 0},
		//{name: "Goroutines", Goroutines: 2},
		//{name: "Goroutines", Goroutines: 4},
		//{name: "Goroutines", Goroutines: 6},
		//{name: "Goroutines", Goroutines: 8},
		//{name: "Goroutines", Goroutines: 10},
	}
	s := newArray(100)

	for _, tt := range tests {
		b.Run(fmt.Sprintf("%s(%d)", tt.name, tt.goroutines), func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = NewSlice(s).Parallel(tt.goroutines).
					Filter(func(v int) bool {
						return v > 100
					}).
					ToSlice()
			}
		})
	}

}

func BenchmarkNativeFilter(b *testing.B) {
	tests := []struct {
		name       string
		goroutines int
		action     func(int, int)
	}{
		{name: "no Parallel", goroutines: 0},
	}
	s := newArray(100)

	for _, tt := range tests {
		b.Run(fmt.Sprintf("%s(%d)", tt.name, tt.goroutines), func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_ = Filter[int](s, func(v int) bool {
					return v > 100
				})
			}
		})
	}

}

func Filter[E any](s []E, predicate func(E) bool) []E {
	clone := make([]E, len(s))
	copy(clone, s)

	newSlice := make([]E, 0, len(s))
	for _, v := range clone {
		if predicate(v) {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

//
//func TestPipelines(t *testing.T) {
//	tests := []struct {
//		name      string
//		input     []string
//		predicate func(v string) bool
//		mapper    func(v string) string
//		want      []string
//	}{
//		{
//			name:      "case",
//			input:     []string{"a", "b", "c"},
//			predicate: func(v string) bool { return v != "b" },
//			mapper: func(v string) string {
//				return v + "1"
//			},
//			want: []string{"a1", "c1"},
//		},
//		{
//			name:      "case",
//			input:     []string{"a", "b"},
//			predicate: func(v string) bool { return v == "c" },
//			want:      []string{},
//		},
//		{
//			name:      "nil",
//			input:     nil,
//			predicate: nil,
//			want:      nil,
//		},
//	}
//	for _, tt := range tests {
//		t.Intermediate(tt.name, func(t *testing.T) {
//			got := NewSlice(tt.input).Filter(tt.predicate).Map(tt.mapper).ToSlice()
//			assert.Equal(t, tt.want, got)
//
//			got = NewSlice(tt.input).Parallel(2).Filter(tt.predicate).Map(tt.mapper).ToSlice()
//			assert.Equal(t, tt.want, got)
//
//			got = NewSliceByComparable(tt.input).Filter(tt.predicate).Map(tt.mapper).ToSlice()
//			assert.Equal(t, tt.want, got)
//
//			got = NewSliceByComparable(tt.input).Parallel(2).Filter(tt.predicate).Map(tt.mapper).ToSlice()
//			assert.Equal(t, tt.want, got)
//
//			got = NewSliceByOrdered(tt.input).Filter(tt.predicate).Map(tt.mapper).ToSlice()
//			assert.Equal(t, tt.want, got)
//
//			got = NewSliceByOrdered(tt.input).Parallel(2).Filter(tt.predicate).Map(tt.mapper).ToSlice()
//			assert.Equal(t, tt.want, got)
//		})
//	}
//}
