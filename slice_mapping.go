package stream

// SliceMappingStream  Need to convert the type of source elements.
// - E elements type
// - MapE map elements type
// - ReduceE reduce elements type
type SliceMappingStream[E any, MapE any, ReduceE any] struct {
	SliceStream[E]
}

// NewSliceByMapping new stream instance, Need to convert the type of source elements.
//
// - E elements type
// - MapE map elements type
// - ReduceE reduce elements type
func NewSliceByMapping[E any, MapE any, ReduceE any](source []E) SliceMappingStream[E, MapE, ReduceE] {
	return SliceMappingStream[E, MapE, ReduceE]{SliceStream: NewSlice(source)}
}

// Map Returns a stream consisting of the results of applying the given function to the elements of this stream.
//
// Support Parallel.
func (stream SliceMappingStream[E, MapE, ReduceE]) Map(mapper func(E) MapE) SliceMappingStream[MapE, MapE, ReduceE] {
	if stream.source == nil {
		return NewSliceByMapping[MapE, MapE, ReduceE](nil)
	}

	termination := func(index int, v E) (isReturn bool, isComplete bool, ret MapE) {
		return true, false, mapper(v)
	}
	return NewSliceByMapping[MapE, MapE, ReduceE](
		pipelineTermination[E, MapE](stream.Pipeline, termination),
	)
}

// Reduce Returns a source consisting of the elements of this stream.
func (stream SliceMappingStream[E, MapE, ReduceE]) Reduce(accumulator func(ReduceE, E) ReduceE) (result ReduceE) {
	stream.evaluation()
	if len(stream.source) == 0 {
		return result
	}

	for _, v := range stream.source {
		result = accumulator(result, v)
	}
	return result
}

// Parallel See: SliceStream.Parallel
func (stream SliceMappingStream[E, MapE, ReduceE]) Parallel(goroutines int) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.Parallel(goroutines)
	return stream
}

// ForEach See: SliceStream.ForEach
func (stream SliceMappingStream[E, MapE, ReduceE]) ForEach(action func(int, E)) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.ForEach(action)
	return stream
}

// Filter See: SliceStream.Filter
func (stream SliceMappingStream[E, MapE, ReduceE]) Filter(predicate func(E) bool) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.Filter(predicate)
	return stream
}

// Limit See: SliceStream.Limit
func (stream SliceMappingStream[E, MapE, ReduceE]) Limit(maxSize int) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.Limit(maxSize)
	return stream
}

// SortFunc See: SliceStream.SortFunc
func (stream SliceMappingStream[E, MapE, ReduceE]) SortFunc(less func(a, b E) bool) SliceMappingStream[E, MapE, ReduceE] {
	stream.SliceStream = stream.SliceStream.SortFunc(less)
	return stream
}
