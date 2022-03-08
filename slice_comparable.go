package stream

import "golang.org/x/exp/slices"

type sliceComparableStream[Elem comparable] struct {
	sliceStream[Elem]
}

func NewSliceByComparable[Elem comparable](v []Elem) sliceComparableStream[Elem] {
	return sliceComparableStream[Elem]{sliceStream: NewSlice(v)}
}

// Distinct Returns a stream consisting of the distinct elements (according to map comparable ) of this stream.
func (stream sliceComparableStream[Elem]) Distinct() sliceComparableStream[Elem] {
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
// according to the slices.EqualFunc
func (stream sliceComparableStream[Elem]) Equal(slice []Elem) bool {
	return slices.Equal(stream.slice, slice)
}

// ForEach Performs an action for each element of this stream.
func (stream sliceComparableStream[Elem]) ForEach(action func(int, Elem)) sliceComparableStream[Elem] {
	stream.sliceStream = stream.sliceStream.ForEach(action)
	return stream
}

// Filter Returns a stream consisting of the elements of this stream that match the given predicate.
func (stream sliceComparableStream[Elem]) Filter(predicate func(Elem) bool) sliceComparableStream[Elem] {
	stream.sliceStream = stream.sliceStream.Filter(predicate)
	return stream
}

// Limit Returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length.
func (stream sliceComparableStream[Elem]) Limit(maxSize int) sliceComparableStream[Elem] {
	stream.sliceStream = stream.sliceStream.Limit(maxSize)
	return stream
}

// Map Returns a stream consisting of the results of applying the given function to the elements of this stream.
func (stream sliceComparableStream[Elem]) Map(mapper func(Elem) Elem) sliceComparableStream[Elem] {
	stream.sliceStream = stream.sliceStream.Map(mapper)
	return stream
}

// SortFunc Returns a stream consisting of the elements of this stream, sorted according to slices.SortFunc.
func (stream sliceComparableStream[Elem]) SortFunc(less func(a, b Elem) bool) sliceComparableStream[Elem] {
	stream.sliceStream = stream.sliceStream.SortFunc(less)
	return stream
}

// SortStableFunc Returns a stream consisting of the elements of this stream, sorted according to slices.SortStableFunc.
func (stream sliceComparableStream[Elem]) SortStableFunc(less func(a, b Elem) bool) sliceComparableStream[Elem] {
	stream.sliceStream = stream.sliceStream.SortStableFunc(less)
	return stream
}
