package stream

import (
	"golang.org/x/exp/slices"
)

type SliceStream[E any] struct {
	slice      []E
	goroutines int
}

// NewSlice new stream instance, generics constraints based on any.
func NewSlice[E any](v []E) SliceStream[E] {
	if v == nil {
		return SliceStream[E]{}
	}
	clone := make([]E, len(v))
	copy(clone, v)
	return SliceStream[E]{slice: clone}
}

// Parallel goroutines > 1 enable parallel, goroutines <= 1 disable parallel
func (stream SliceStream[E]) Parallel(goroutines int) SliceStream[E] {
	stream.goroutines = goroutines
	return stream
}

// IsParallel Whether parallelism is enabled
func (stream SliceStream[E]) IsParallel() bool {
	if stream.goroutines > 1 {
		return true
	}
	return false
}

// At Returns the element at the given index. Accepts negative integers, which count back from the last item.
func (stream SliceStream[E]) At(index int) (elem E) {
	l := len(stream.slice)
	if index < 0 {
		index = index + l
	}
	if l == 0 || index < 0 || index >= l {
		return
	}
	return stream.slice[index]
}

// AllMatch Returns whether all elements in the stream match the provided predicate.
// If the slice is empty or nil then true is returned.
//
// Support Parallel.
func (stream SliceStream[E]) AllMatch(predicate func(E) bool) bool {
	if stream.IsParallel() {
		handler := func(index int, v E) (isReturn bool, result bool) {
			return !predicate(v), false
		}

		results := ParallelProcess[ParallelFirst[E, bool], E, bool](
			stream.goroutines,
			stream.slice,
			handler,
		)

		if len(results) > 0 {
			return results[0]
		}
		return true
	}

	for _, v := range stream.slice {
		if !predicate(v) {
			return false
		}
	}
	return true

}

// AnyMatch Returns whether any elements in the stream match the provided predicate.
// If the slice is empty or nil then false is returned.
//
// Support Parallel.
func (stream SliceStream[E]) AnyMatch(predicate func(E) bool) bool {
	if stream.IsParallel() {
		handler := func(index int, v E) (isReturn bool, result bool) {
			return predicate(v), true
		}
		results := ParallelProcess[ParallelFirst[E, bool], E, bool](
			stream.goroutines,
			stream.slice,
			handler,
		)

		if len(results) > 0 {
			return results[0]
		}
		return false
	}

	for _, v := range stream.slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

// Append appends elements to the end of this stream
func (stream SliceStream[E]) Append(elements ...E) SliceStream[E] {
	newSlice := make([]E, 0, len(stream.slice)+len(elements))
	newSlice = append(newSlice, stream.slice...)
	newSlice = append(newSlice, elements...)
	stream.slice = newSlice
	return stream
}

// Count Returns the count of elements in this stream.
func (stream SliceStream[E]) Count() int {
	return len(stream.slice)
}

// EqualFunc Returns whether the slice in the stream is equal to the destination slice.
// Equal according to the slices.EqualFunc
func (stream SliceStream[E]) EqualFunc(dest []E, equal func(E, E) bool) bool {
	return slices.EqualFunc(stream.slice, dest, equal)
}

// ForEach Performs an action for each element of this stream.
//
// Support Parallel.
// Parallel side effects are not executed in the original order of stream elements.
func (stream SliceStream[E]) ForEach(action func(int, E)) SliceStream[E] {
	if stream.slice == nil {
		return stream
	}

	if stream.IsParallel() {
		handler := func(index int, v E) (isReturn bool, result E) {
			action(index, v)
			return false, result
		}
		ParallelProcess[ParallelAction[E, E], E, E](
			stream.goroutines,
			stream.slice,
			handler)
		return stream
	}

	for i, v := range stream.slice {
		action(i, v)
	}
	return stream
}

// First Returns the first element in the stream.
// If the slice is empty or nil then E Type default value is returned.
func (stream SliceStream[E]) First() (elem E) {
	if len(stream.slice) == 0 {
		return
	}
	return stream.slice[0]
}

// FindFunc Returns the index of the first element in the stream that matches the provided predicate.
// If not found then -1 is returned.
//
// Support Parallel.
// Parallel side effect is that the element found may not be the first to appear
func (stream SliceStream[E]) FindFunc(predicate func(E) bool) int {
	if stream.IsParallel() {
		handler := func(index int, v E) (isReturn bool, result int) {
			return predicate(v), index
		}
		results := ParallelProcess[ParallelFirst[E, int], E, int](
			stream.goroutines,
			stream.slice,
			handler)

		if len(results) > 0 {
			return results[0]
		}
		return -1
	}

	for i, v := range stream.slice {
		if predicate(v) {
			return i
		}
	}
	return -1
}

// Filter Returns a stream consisting of the elements of this stream that match the given predicate.
//
// Support Parallel.
func (stream SliceStream[E]) Filter(predicate func(E) bool) SliceStream[E] {
	if stream.slice == nil {
		return stream
	}

	if stream.IsParallel() {
		handler := func(index int, v E) (isReturn bool, result E) {
			return predicate(v), v
		}
		stream.slice = ParallelProcess[ParallelAll[E, E], E, E](
			stream.goroutines,
			stream.slice,
			handler)
		return stream
	}

	newSlice := make([]E, 0)
	for _, v := range stream.slice {
		if predicate(v) {
			newSlice = append(newSlice, v)
		}
	}
	stream.slice = newSlice
	return stream
}

// Insert inserts the values v... into s at index
// If index is out of range then use Append to the end
func (stream SliceStream[E]) Insert(index int, elements ...E) SliceStream[E] {
	if len(stream.slice) <= index {
		return stream.Append(elements...)
	}
	stream.slice = slices.Insert(stream.slice, index, elements...)
	return stream
}

// Delete Removes the elements s[i:j] from this stream, returning the modified stream.
// If the slice is empty or nil then do nothing
func (stream SliceStream[E]) Delete(i, j int) SliceStream[E] {
	if len(stream.slice) == 0 {
		return stream
	}
	stream.slice = slices.Delete(stream.slice, i, j)
	return stream
}

// IsSortedFunc Returns whether stream is sorted in ascending order.
// Compare according to the less function
// - less: return a > b
// If the slice is empty or nil then true is returned.
func (stream SliceStream[E]) IsSortedFunc(less func(a, b E) bool) bool {
	return slices.IsSortedFunc(stream.slice, less)
}

// Limit Returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length.
func (stream SliceStream[E]) Limit(maxSize int) SliceStream[E] {
	if stream.slice == nil {
		return stream
	}

	newSlice := make([]E, 0)
	for i := 0; i < len(stream.slice) && i < maxSize; i++ {
		newSlice = append(newSlice, stream.slice[i])
	}
	stream.slice = newSlice
	return stream
}

// Map Returns a stream consisting of the results of applying the given function to the elements of this stream.
//
// Support Parallel.
func (stream SliceStream[E]) Map(mapper func(E) E) SliceStream[E] {
	if stream.slice == nil {
		return stream
	}

	if stream.IsParallel() {
		handler := func(index int, v E) (isReturn bool, result E) {
			return true, mapper(v)
		}
		stream.slice = ParallelProcess[ParallelAll[E, E], E, E](
			stream.goroutines,
			stream.slice,
			handler)
		return stream
	}

	for i, v := range stream.slice {
		stream.slice[i] = mapper(v)
	}
	return stream
}

// MaxFunc Returns the maximum element of this stream.
// - less: return a > b
// If the slice is empty or nil then E Type default value is returned.
func (stream SliceStream[E]) MaxFunc(less func(a, b E) bool) (max E) {
	for i, v := range stream.slice {
		if less(v, max) || i == 0 {
			max = v
		}
	}
	return max
}

// MinFunc Returns the minimum element of this stream.
// - less: return a < b
// If the slice is empty or nil then E Type default value is returned.
func (stream SliceOrderedStream[E]) MinFunc(less func(a, b E) bool) (min E) {
	for i, v := range stream.slice {
		if less(v, min) || i == 0 {
			min = v
		}
	}
	return min
}

// Reduce Returns a slice consisting of the elements of this stream.
func (stream SliceStream[E]) Reduce(accumulator func(E, E) E) (result E) {
	if len(stream.slice) == 0 {
		return result
	}

	for _, v := range stream.slice {
		result = accumulator(result, v)
	}
	return result
}

// SortFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortFunc.
func (stream SliceStream[E]) SortFunc(less func(a, b E) bool) SliceStream[E] {
	slices.SortFunc(stream.slice, less)
	return stream
}

// SortStableFunc Returns a sorted stream consisting of the elements of this stream.
// Sorted according to slices.SortStableFunc.
func (stream SliceStream[E]) SortStableFunc(less func(a, b E) bool) SliceStream[E] {
	slices.SortStableFunc(stream.slice, less)
	return stream
}

// ToSlice Returns a slice in the stream
func (stream SliceStream[E]) ToSlice() []E {
	return stream.slice
}
