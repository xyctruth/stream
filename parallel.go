package stream

// ParallelHandlerFunc Parallel processing Handler
//
// - E elements type
// - R result type
type ParallelHandlerFunc[E any, R any] func(index int, elem E) (isReturn bool, result R)

// Parallel Provides parallel processing capabilities.
//
// - E elements type
// - R result type
//
// ParallelFirst parallel processing ends as soon as the first return value is obtained.
//
// ParallelAll All elements need to be processed in parallel, all return values are obtained, and then the parallel is ended.
//
// ParallelAction All elements need to be processed in parallel, no return value required
type Parallel[E any, R any] interface {
	Process(goroutines int, slice []E, handler ParallelHandlerFunc[E, R]) []R
}

func ParallelProcess[T Parallel[E, R], E any, R any](
	goroutines int,
	slice []E,
	handler ParallelHandlerFunc[E, R]) []R {
	var p T
	return p.Process(goroutines, slice, handler)
}

// part  Uniform slices
// This selects a half-open range which includes the first element, but excludes the last.
type part struct {
	low  int //includes index
	high int //excludes index
}

// partition Given a specified slice, evenly part according to the slice.
func partition[E any](slice []E, goroutines int) []part {
	l := len(slice)
	if l == 0 {
		return nil
	}
	if goroutines > l {
		goroutines = l
	}
	partitions := make([]part, 0, goroutines)
	size := l / goroutines
	rem := l % goroutines
	low := 0
	high := size
	for i := 0; i < goroutines; i++ {
		partitions = append(partitions, part{low, high})
		low = high
		high = high + size
		if i+rem+1 >= goroutines {
			high++
		}
	}
	return partitions
}
