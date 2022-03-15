package stream

import "golang.org/x/exp/slices"

type SliceComparableStream[Elem comparable] struct {
	SliceStream[Elem]
}

// NewSliceByComparable  new stream instance, generics constraints based on comparable
func NewSliceByComparable[Elem comparable](v []Elem) SliceComparableStream[Elem] {
	return SliceComparableStream[Elem]{SliceStream: NewSlice(v)}
}

// Parallel goroutines > 1 enable All, goroutines <= 1 disable All
func (stream SliceComparableStream[Elem]) Parallel(goroutines int) SliceComparableStream[Elem] {
	stream.SliceStream = stream.SliceStream.Parallel(goroutines)
	return stream
}

// Distinct Returns a stream consisting of the distinct elements of this stream.
// Remove duplicate according to map comparable.
func (stream SliceComparableStream[Elem]) Distinct() SliceComparableStream[Elem] {
	if stream.slice == nil {
		return stream
	}

	newSlice := make([]Elem, 0)
	distinct := map[Elem]struct{}{}
	for _, v := range stream.slice {
		if _, ok := distinct[v]; ok {
			continue
		}
		distinct[v] = struct{}{}
		newSlice = append(newSlice, v)
	}
	stream.slice = newSlice
	return stream
}

// Equal Returns whether the slice in the stream is equal to the destination slice.
// Equal according to the slices.Equal.
func (stream SliceComparableStream[Elem]) Equal(dest []Elem) bool {
	return slices.Equal(stream.slice, dest)
}

// Find Returns the index of the first element in the stream that matches the target element.
// If not found then -1 is returned.
func (stream SliceComparableStream[Elem]) Find(dest Elem) int {
	for i, v := range stream.slice {
		if v == dest {
			return i
		}
	}
	return -1
}

// ForEach Performs an action for each element of this stream.
func (stream SliceComparableStream[Elem]) ForEach(action func(int, Elem)) SliceComparableStream[Elem] {
	stream.SliceStream = stream.SliceStream.ForEach(action)
	return stream
}

// Filter Returns a stream consisting of the elements of this stream that match the given predicate.
func (stream SliceComparableStream[Elem]) Filter(predicate func(Elem) bool) SliceComparableStream[Elem] {
	stream.SliceStream = stream.SliceStream.Filter(predicate)
	return stream
}

// Limit Returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length.
func (stream SliceComparableStream[Elem]) Limit(maxSize int) SliceComparableStream[Elem] {
	stream.SliceStream = stream.SliceStream.Limit(maxSize)
	return stream
}

// Map Returns a stream consisting of the results of applying the given function to the elements of this stream.
func (stream SliceComparableStream[Elem]) Map(mapper func(Elem) Elem) SliceComparableStream[Elem] {
	stream.SliceStream = stream.SliceStream.Map(mapper)
	return stream
}

// SortFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortFunc.
func (stream SliceComparableStream[Elem]) SortFunc(less func(a, b Elem) bool) SliceComparableStream[Elem] {
	stream.SliceStream = stream.SliceStream.SortFunc(less)
	return stream
}

// SortStableFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortStableFunc.
func (stream SliceComparableStream[Elem]) SortStableFunc(less func(a, b Elem) bool) SliceComparableStream[Elem] {
	stream.SliceStream = stream.SliceStream.SortStableFunc(less)
	return stream
}
