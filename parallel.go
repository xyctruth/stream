package stream

func Parallel[E any, R any](
	goroutines int,
	slice []E,
	handler func(index int, elem E) (isReturn bool, isComplete bool, result R)) []R {

	partitions := partition(slice, goroutines)
	resultChs := make([]chan []R, len(partitions))

	for i, pa := range partitions {
		resultChs[i] = make(chan []R)
		go parallelDo(slice, handler, resultChs[i], pa)
	}

	result := parallelResulted(resultChs, len(slice))
	return result
}

func parallelDo[E any, R any](
	slice []E,
	handler func(index int, elem E) (isReturn bool, isComplete bool, result R),
	resultCh chan []R,
	pa part) {

	defer close(resultCh)
	ret := make([]R, 0, pa.high-pa.low)

	for i := pa.low; i < pa.high; i++ {
		isReturn, isComplete, r := handler(i, slice[i])
		if !isReturn {
			continue
		}
		ret = append(ret, r)
		if isComplete {
			break
		}
	}

	if len(ret) > 0 {
		resultCh <- ret
	}
	return
}

func parallelResulted[R any](resultChs []chan []R, cap int) []R {
	results := make([]R, 0, cap)
	for _, resultCh := range resultChs {
		for result := range resultCh {
			results = append(results, result...)
		}
	}
	return results
}

// part  Uniform slices
// This selects a half-open range which includes the first element, but excludes the last.
type part struct {
	low  int //includes index
	high int //excludes index
}

// partition Given a specified source, evenly part according to the source.
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
