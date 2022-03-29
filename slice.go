package stream

import (
	"golang.org/x/exp/slices"
)

// SliceStream Generics constraints based on any
type SliceStream[E any] struct {
	*Pipe[E]
}

// NewSlice new stream instance, generics constraints based on any.
func NewSlice[E any](source []E) SliceStream[E] {
	return SliceStream[E]{Pipe: &Pipe[E]{source: source}}
}

// Parallel Goroutines > 1 enable parallel, Goroutines <= 1 disable parallel
func (stream SliceStream[E]) Parallel(goroutines int) SliceStream[E] {
	stream.goroutines = goroutines
	return stream
}

// At Returns the element at the given index. Accepts negative integers, which count back from the last item.
func (stream SliceStream[E]) At(index int) (elem E) {
	stream.evaluation()
	l := len(stream.source)
	if index < 0 {
		index = index + l
	}
	if l == 0 || index < 0 || index >= l {
		return
	}
	return stream.source[index]
}

// AllMatch Returns whether all elements in the stream match the provided predicate.
// If the source is empty or nil then true is returned.
//
// Support Parallel.
func (stream SliceStream[E]) AllMatch(predicate func(E) bool) bool {
	termination := func(index int, v E) (isReturn bool, isComplete bool, ret bool) {
		return !predicate(v), true, false
	}
	results := stream.evaluationBoolTermination(termination)
	if len(results) > 0 {
		return results[0]
	}
	return true
}

// AnyMatch Returns whether any elements in the stream match the provided predicate.
// If the source is empty or nil then false is returned.
//
// Support Parallel.
func (stream SliceStream[E]) AnyMatch(predicate func(E) bool) bool {
	termination := func(index int, v E) (isReturn bool, isComplete bool, ret bool) {
		return predicate(v), true, true
	}
	results := stream.evaluationBoolTermination(termination)
	if len(results) > 0 {
		return results[0]
	}
	return false
}

// Append appends elements to the end of this stream
func (stream SliceStream[E]) Append(elements ...E) SliceStream[E] {
	stream.evaluation()
	newSlice := make([]E, 0, len(stream.source)+len(elements))
	newSlice = append(newSlice, stream.source...)
	newSlice = append(newSlice, elements...)
	stream.source = newSlice
	return stream
}

// Count Returns the count of elements in this stream.
func (stream SliceStream[E]) Count() int {
	stream.evaluation()
	return len(stream.source)
}

// EqualFunc Returns whether the source in the stream is equal to the destination source.
// Equal according to the slices.EqualFunc
func (stream SliceStream[E]) EqualFunc(dest []E, equal func(E, E) bool) bool {
	stream.evaluation()
	return slices.EqualFunc(stream.source, dest, equal)
}

// ForEach Performs an action for each element of this stream.
//
// Support Parallel.
// Parallel side effects are not executed in the original order of stream elements.
func (stream SliceStream[E]) ForEach(action func(int, E)) SliceStream[E] {
	if stream.source == nil {
		return stream
	}
	stage := func(index int, v E) (isReturn bool, isComplete bool, result E) {
		action(index, v)
		return true, false, v
	}
	stream.Pipe.AddStage(stage)
	stream.evaluation()
	return stream
}

// First Returns the first element in the stream.
// If the source is empty or nil then E Type default value is returned.
func (stream SliceStream[E]) First() (elem E) {
	stream.evaluation()
	if len(stream.source) == 0 {
		return
	}
	return stream.source[0]
}

// FindFunc Returns the index of the first element in the stream that matches the provided predicate.
// If not found then -1 is returned.
//
// Support Parallel.
// Parallel side effect is that the element found may not be the first to appear
func (stream SliceStream[E]) FindFunc(predicate func(E) bool) int {
	stage := func(index int, v E) (isReturn bool, isComplete bool, ret int) {
		return predicate(v), true, index
	}
	results := stream.evaluationIntTermination(stage)
	if len(results) > 0 {
		return results[0]
	}
	return -1
}

// Filter Returns a stream consisting of the elements of this stream that match the given predicate.
//
// Support Parallel.
func (stream SliceStream[E]) Filter(predicate func(E) bool) SliceStream[E] {
	stage := func(index int, e E) (isReturn bool, isComplete bool, ret E) {
		return predicate(e), false, e
	}
	stream.Pipe.AddStage(stage)
	return stream
}

// Insert inserts the values source... into s at index
// If index is out of range then use Append to the end
func (stream SliceStream[E]) Insert(index int, elements ...E) SliceStream[E] {
	stream.evaluation()
	if len(stream.source) <= index {
		return stream.Append(elements...)
	}
	stream.source = slices.Insert(stream.source, index, elements...)
	return stream
}

// Delete Removes the elements s[i:j] from this stream, returning the modified stream.
// If the source is empty or nil then do nothing
func (stream SliceStream[E]) Delete(i, j int) SliceStream[E] {
	stream.evaluation()
	if len(stream.source) == 0 {
		return stream
	}
	stream.source = slices.Delete(stream.source, i, j)
	return stream
}

// IsSortedFunc Returns whether stream is sorted in ascending order.
// Compare according to the less function
// - less: return a > b
// If the source is empty or nil then true is returned.
func (stream SliceStream[E]) IsSortedFunc(less func(a, b E) bool) bool {
	stream.evaluation()
	return slices.IsSortedFunc(stream.source, less)
}

// Limit Returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length.
func (stream SliceStream[E]) Limit(maxSize int) SliceStream[E] {
	stream.evaluation()
	if stream.source == nil {
		return stream
	}
	newSlice := make([]E, 0, maxSize)
	for i := 0; i < len(stream.source) && i < maxSize; i++ {
		newSlice = append(newSlice, stream.source[i])
	}
	stream.source = newSlice
	return stream
}

// Map Returns a stream consisting of the results of applying the given function to the elements of this stream.
//
// Support Parallel.
func (stream SliceStream[E]) Map(mapper func(E) E) SliceStream[E] {
	stage := func(index int, v E) (isReturn bool, isComplete bool, ret E) {
		return true, false, mapper(v)
	}
	stream.Pipe.AddStage(stage)
	return stream
}

// MaxFunc Returns the maximum element of this stream.
// - less: return a > b
// If the source is empty or nil then E Type default value is returned.
func (stream SliceStream[E]) MaxFunc(less func(a, b E) bool) (max E) {
	stream.evaluation()
	for i, v := range stream.source {
		if less(v, max) || i == 0 {
			max = v
		}
	}
	return max
}

// MinFunc Returns the minimum element of this stream.
// - less: return a < b
// If the source is empty or nil then E Type default value is returned.
func (stream SliceOrderedStream[E]) MinFunc(less func(a, b E) bool) (min E) {
	stream.evaluation()
	for i, v := range stream.source {
		if less(v, min) || i == 0 {
			min = v
		}
	}
	return min
}

// Reduce Returns a source consisting of the elements of this stream.
func (stream SliceStream[E]) Reduce(accumulator func(E, E) E) (result E) {
	stream.evaluation()
	if len(stream.source) == 0 {
		return result
	}

	for _, v := range stream.source {
		result = accumulator(result, v)
	}
	return result
}

// SortFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortFunc.
func (stream SliceStream[E]) SortFunc(less func(a, b E) bool) SliceStream[E] {
	stream.evaluation()
	slices.SortFunc(stream.source, less)
	return stream
}

// SortStableFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortStableFunc.
func (stream SliceStream[E]) SortStableFunc(less func(a, b E) bool) SliceStream[E] {
	stream.evaluation()
	slices.SortStableFunc(stream.source, less)
	return stream
}

// ToSlice Returns a source in the stream
func (stream SliceStream[E]) ToSlice() []E {
	stream.evaluation()
	return stream.source
}

func (stream *SliceStream[E]) evaluation() {
	stream.Pipe.Run()
}

func (stream *SliceStream[E]) evaluationBoolTermination(termination Stage[E, bool]) (ret []bool) {
	ret = PipeByTermination[E, bool](stream.Pipe, termination)
	stream.stages = nil
	return ret
}

func (stream *SliceStream[E]) evaluationIntTermination(termination Stage[E, int]) (ret []int) {
	ret = PipeByTermination[E, int](stream.Pipe, termination)
	stream.stages = nil
	return ret
}
