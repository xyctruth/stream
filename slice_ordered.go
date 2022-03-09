package stream

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type sliceOrderedStream[Elem constraints.Ordered] struct {
	sliceComparableStream[Elem]
}

// NewSliceByOrdered new stream instance, generics constraints based on constraints.Ordered
func NewSliceByOrdered[Elem constraints.Ordered](v []Elem) sliceOrderedStream[Elem] {
	return sliceOrderedStream[Elem]{sliceComparableStream: NewSliceByComparable(v)}
}

// Parallel cores > 1 enable parallel, cores <= 1 disable parallel
func (stream sliceOrderedStream[Elem]) Parallel(cores int) sliceOrderedStream[Elem] {
	stream.sliceComparableStream = stream.sliceComparableStream.Parallel(cores)
	return stream
}

// Max Returns the maximum element of this stream.
// Compare according to the constraints.Ordered.
// If the slice is empty or nil then Elem Type default value is returned.
func (stream sliceOrderedStream[Elem]) Max() Elem {
	var max Elem
	for _, v := range stream.slice {
		if v > max {
			max = v
		}
	}
	return max
}

// Min Returns the minimum element of this stream.
// Compare according to the constraints.Ordered.
// If the slice is empty or nil then Elem Type default value is returned.
func (stream sliceOrderedStream[Elem]) Min() Elem {
	var min Elem
	for i, v := range stream.slice {
		if v < min || i == 0 {
			min = v
		}
	}
	return min
}

// Sort Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.Sort.
func (stream sliceOrderedStream[Elem]) Sort() sliceOrderedStream[Elem] {
	slices.Sort(stream.slice)
	return stream
}

// Distinct Returns a stream consisting of the distinct elements of this stream.
// Remove duplicate according to map comparable.
func (stream sliceOrderedStream[Elem]) Distinct() sliceOrderedStream[Elem] {
	stream.sliceComparableStream = stream.sliceComparableStream.Distinct()
	return stream
}

// ForEach Performs an action for each element of this stream.
func (stream sliceOrderedStream[Elem]) ForEach(action func(int, Elem)) sliceOrderedStream[Elem] {
	stream.sliceComparableStream = stream.sliceComparableStream.ForEach(action)
	return stream
}

// Filter Returns a stream consisting of the elements of this stream that match the given predicate.
func (stream sliceOrderedStream[Elem]) Filter(predicate func(Elem) bool) sliceOrderedStream[Elem] {
	stream.sliceComparableStream = stream.sliceComparableStream.Filter(predicate)
	return stream
}

// Limit Returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length.
func (stream sliceOrderedStream[Elem]) Limit(maxSize int) sliceOrderedStream[Elem] {
	stream.sliceComparableStream = stream.sliceComparableStream.Limit(maxSize)
	return stream
}

// Map Returns a stream consisting of the results of applying the given function to the elements of this stream.
func (stream sliceOrderedStream[Elem]) Map(mapper func(Elem) Elem) sliceOrderedStream[Elem] {
	stream.sliceComparableStream = stream.sliceComparableStream.Map(mapper)
	return stream
}

// SortFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortFunc.
func (stream sliceOrderedStream[Elem]) SortFunc(less func(a, b Elem) bool) sliceOrderedStream[Elem] {
	stream.sliceComparableStream = stream.sliceComparableStream.SortFunc(less)
	return stream
}

// SortStableFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortStableFunc.
func (stream sliceOrderedStream[Elem]) SortStableFunc(less func(a, b Elem) bool) sliceOrderedStream[Elem] {
	stream.sliceComparableStream = stream.sliceComparableStream.SortStableFunc(less)
	return stream
}
