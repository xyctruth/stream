package stream

// SliceMappingStream  Need to convert the type of slice elements.
// - E elements type
// - MapE map elements type
// - ReduceE reduce elements type
type SliceMappingStream[E any, MapE any, ReduceE any] struct {
	SliceStream[E]
}

// NewSliceByMapping new stream instance, Need to convert the type of slice elements.
//
// - E elements type
// - MapE map elements type
// - ReduceE reduce elements type
func NewSliceByMapping[E any, MapE any, ReduceE any](v []E) SliceMappingStream[E, MapE, ReduceE] {
	return SliceMappingStream[E, MapE, ReduceE]{SliceStream: NewSlice(v)}
}

// Map Returns a stream consisting of the results of applying the given function to the elements of this stream.
//
// Support Parallel.
func (stream SliceMappingStream[E, MapE, ReduceE]) Map(mapper func(E) MapE) SliceMappingStream[MapE, MapE, ReduceE] {
	if stream.slice == nil {
		return NewSliceByMapping[MapE, MapE, ReduceE](nil)
	}
	if stream.IsParallel() {
		handler := func(index int, v E) (isReturn bool, result MapE) {
			return true, mapper(v)
		}
		newSlice := ParallelProcess[ParallelAll[E, MapE], E, MapE](
			stream.goroutines,
			stream.slice,
			handler)
		return NewSliceByMapping[MapE, MapE, ReduceE](newSlice)
	}

	newSlice := make([]MapE, len(stream.slice))
	for i, v := range stream.slice {
		newSlice[i] = mapper(v)
	}
	return NewSliceByMapping[MapE, MapE, ReduceE](newSlice)
}

// Reduce Returns a slice consisting of the elements of this stream.
func (stream SliceMappingStream[E, MapE, ReduceE]) Reduce(accumulator func(ReduceE, E) ReduceE) (result ReduceE) {
	if len(stream.slice) == 0 {
		return result
	}

	for _, v := range stream.slice {
		result = accumulator(result, v)
	}
	return result
}

// Parallel goroutines > 1 enable All, goroutines <= 1 disable All
func (stream SliceMappingStream[E, MapE, ReduceE]) Parallel(goroutines int) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.Parallel(goroutines)
	return stream
}

// ForEach Performs an action for each element of this stream.
func (stream SliceMappingStream[E, MapE, ReduceE]) ForEach(action func(int, E)) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.ForEach(action)
	return stream
}

// Filter Returns a stream consisting of the elements of this stream that match the given predicate.
func (stream SliceMappingStream[E, MapE, ReduceE]) Filter(predicate func(E) bool) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.Filter(predicate)
	return stream
}

// Limit Returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length.
func (stream SliceMappingStream[E, MapE, ReduceE]) Limit(maxSize int) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.Limit(maxSize)
	return stream
}

// SortFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortFunc.
func (stream SliceMappingStream[E, MapE, ReduceE]) SortFunc(less func(a, b E) bool) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.SortFunc(less)
	return stream
}

// SortStableFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortStableFunc.
func (stream SliceMappingStream[E, MapE, ReduceE]) SortStableFunc(less func(a, b E) bool) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.SortStableFunc(less)
	return stream
}
