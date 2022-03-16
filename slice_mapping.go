package stream

type SliceMappingStream[Elem any, MapElem any, ReduceElem any] struct {
	SliceStream[Elem]
}

// NewSliceByMapping new stream instance, Need to convert the type of slice elements.
func NewSliceByMapping[Elem any, MapElem any, ReduceElem any](v []Elem) SliceMappingStream[Elem, MapElem, ReduceElem] {
	return SliceMappingStream[Elem, MapElem, ReduceElem]{SliceStream: NewSlice(v)}
}

// Map Returns a stream consisting of the results of applying the given function to the elements of this stream.
//
// Support Parallel.
func (stream SliceMappingStream[Elem, MapElem, ReduceElem]) Map(mapper func(Elem) MapElem) SliceMappingStream[MapElem, MapElem, ReduceElem] {
	if stream.slice == nil {
		return NewSliceByMapping[MapElem, MapElem, ReduceElem](nil)
	}
	if stream.IsParallel() {
		handler := func(index int, v Elem) (isReturn bool, taskResult MapElem) {
			return true, mapper(v)
		}
		newSlice := ParallelProcess[ParallelAll[Elem, MapElem], Elem, MapElem](
			stream.goroutines,
			stream.slice,
			handler)
		return NewSliceByMapping[MapElem, MapElem, ReduceElem](newSlice)
	}

	newSlice := make([]MapElem, len(stream.slice))
	for i, v := range stream.slice {
		newSlice[i] = mapper(v)
	}
	return NewSliceByMapping[MapElem, MapElem, ReduceElem](newSlice)
}

// Reduce Returns a slice consisting of the elements of this stream.
func (stream SliceMappingStream[Elem, MapElem, ReduceElem]) Reduce(accumulator func(ReduceElem, Elem) ReduceElem) (result ReduceElem) {
	if len(stream.slice) == 0 {
		return result
	}

	for _, v := range stream.slice {
		result = accumulator(result, v)
	}
	return result
}

// Parallel goroutines > 1 enable All, goroutines <= 1 disable All
func (stream SliceMappingStream[Elem, MapElem, ReduceElem]) Parallel(goroutines int) SliceMappingStream[Elem, MapElem, ReduceElem] {
	stream.SliceStream = stream.SliceStream.Parallel(goroutines)
	return stream
}

// ForEach Performs an action for each element of this stream.
func (stream SliceMappingStream[Elem, MapElem, ReduceElem]) ForEach(action func(int, Elem)) SliceMappingStream[Elem, MapElem, ReduceElem] {
	stream.SliceStream = stream.SliceStream.ForEach(action)
	return stream
}

// Filter Returns a stream consisting of the elements of this stream that match the given predicate.
func (stream SliceMappingStream[Elem, MapElem, ReduceElem]) Filter(predicate func(Elem) bool) SliceMappingStream[Elem, MapElem, ReduceElem] {
	stream.SliceStream = stream.SliceStream.Filter(predicate)
	return stream
}

// Limit Returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length.
func (stream SliceMappingStream[Elem, MapElem, ReduceElem]) Limit(maxSize int) SliceMappingStream[Elem, MapElem, ReduceElem] {
	stream.SliceStream = stream.SliceStream.Limit(maxSize)
	return stream
}

// SortFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortFunc.
func (stream SliceMappingStream[Elem, MapElem, ReduceElem]) SortFunc(less func(a, b Elem) bool) SliceMappingStream[Elem, MapElem, ReduceElem] {
	stream.SliceStream = stream.SliceStream.SortFunc(less)
	return stream
}

// SortStableFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortStableFunc.
func (stream SliceMappingStream[Elem, MapElem, ReduceElem]) SortStableFunc(less func(a, b Elem) bool) SliceMappingStream[Elem, MapElem, ReduceElem] {
	stream.SliceStream = stream.SliceStream.SortStableFunc(less)
	return stream
}
