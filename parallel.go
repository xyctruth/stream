package stream

type ParallelHandler[Elem any, TaskResult any] func(index int, elem Elem) (isReturn bool, taskResult TaskResult)

// Parallel Provides parallel processing capabilities.
//
// ParallelFirst parallel processing ends as soon as the first return value is obtained.
// ParallelAll All elements need to be processed in parallel, all return values are obtained, and then the parallel is ended.
// ParallelAction All elements need to be processed in parallel, no return value required
type Parallel[Elem any, Result any] interface {
	Process(goroutines int, slice []Elem, handler ParallelHandler[Elem, Result]) []Result
}

func ParallelProcess[T Parallel[Elem, Result], Elem any, Result any](
	goroutines int,
	slice []Elem,
	handler ParallelHandler[Elem, Result]) []Result {
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
func partition[Elem any](slice []Elem, goroutines int) []part {
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
