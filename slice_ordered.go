package stream

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type SliceOrderedStream[E constraints.Ordered] struct {
	SliceComparableStream[E]
}

// NewSliceByOrdered new stream instance, generics constraints based on constraints.Ordered
func NewSliceByOrdered[E constraints.Ordered](v []E) SliceOrderedStream[E] {
	return SliceOrderedStream[E]{SliceComparableStream: NewSliceByComparable(v)}
}

// Parallel goroutines > 1 enable All, goroutines <= 1 disable All
func (stream SliceOrderedStream[E]) Parallel(goroutines int) SliceOrderedStream[E] {
	stream.SliceComparableStream = stream.SliceComparableStream.Parallel(goroutines)
	return stream
}

// IsSorted reports whether x is sorted in ascending order.
// Compare according to the constraints.Ordered.
// If the slice is empty or nil then true is returned.
func (stream SliceOrderedStream[E]) IsSorted() bool {
	return slices.IsSorted(stream.slice)
}

// Max Returns the maximum element of this stream.
// Compare according to the constraints.Ordered.
// If the slice is empty or nil then E Type default value is returned.
func (stream SliceOrderedStream[E]) Max() E {
	var max E
	for i, v := range stream.slice {
		if v > max || i == 0 {
			max = v
		}
	}
	return max
}

// Min Returns the minimum element of this stream.
// Compare according to the constraints.Ordered.
// If the slice is empty or nil then E Type default value is returned.
func (stream SliceOrderedStream[E]) Min() E {
	var min E
	for i, v := range stream.slice {
		if v < min || i == 0 {
			min = v
		}
	}
	return min
}

// Sort Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.Sort.
func (stream SliceOrderedStream[E]) Sort() SliceOrderedStream[E] {
	slices.Sort(stream.slice)
	return stream
}

// Distinct Returns a stream consisting of the distinct elements of this stream.
// Remove duplicate according to map comparable.
func (stream SliceOrderedStream[E]) Distinct() SliceOrderedStream[E] {
	stream.SliceComparableStream = stream.SliceComparableStream.Distinct()
	return stream
}

// ForEach Performs an action for each element of this stream.
func (stream SliceOrderedStream[E]) ForEach(action func(int, E)) SliceOrderedStream[E] {
	stream.SliceStream = stream.SliceStream.ForEach(action)
	return stream
}

// Filter Returns a stream consisting of the elements of this stream that match the given predicate.
func (stream SliceOrderedStream[E]) Filter(predicate func(E) bool) SliceOrderedStream[E] {
	stream.SliceStream = stream.SliceStream.Filter(predicate)
	return stream
}

// Limit Returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length.
func (stream SliceOrderedStream[E]) Limit(maxSize int) SliceOrderedStream[E] {
	stream.SliceStream = stream.SliceStream.Limit(maxSize)
	return stream
}

// Map Returns a stream consisting of the results of applying the given function to the elements of this stream.
func (stream SliceOrderedStream[E]) Map(mapper func(E) E) SliceOrderedStream[E] {
	stream.SliceStream = stream.SliceStream.Map(mapper)
	return stream
}

// SortFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortFunc.
func (stream SliceOrderedStream[E]) SortFunc(less func(a, b E) bool) SliceOrderedStream[E] {
	stream.SliceStream = stream.SliceStream.SortFunc(less)
	return stream
}

// SortStableFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortStableFunc.
func (stream SliceOrderedStream[E]) SortStableFunc(less func(a, b E) bool) SliceOrderedStream[E] {
	stream.SliceStream = stream.SliceStream.SortStableFunc(less)
	return stream
}
