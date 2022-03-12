package stream

func singleResultHandler[Elem any](defaultVal Elem) parallelResultHandler[Elem, Elem] {
	return func(taskResultChs []chan []Elem) Elem {
		for _, ch := range taskResultChs {
			for result := range ch {
				for _, r := range result {
					return r
				}
			}
		}
		return defaultVal
	}
}

func multipleResultHandler[Elem any](count int) parallelResultHandler[Elem, []Elem] {
	return func(taskResultChs []chan []Elem) []Elem {
		results := make([]Elem, 0, count)
		for _, ch := range taskResultChs {
			for result := range ch {
				for _, r := range result {
					results = append(results, r)
				}
			}
		}
		return results
	}
}

// partition  Uniform slices
// This selects a half-open range which includes the first element, but excludes the last.
type partition struct {
	low  int //includes index
	high int //excludes index
}

func partitionHandler[Elem any](slice []Elem, goroutines int) []partition {
	l := len(slice)
	if goroutines > l {
		goroutines = l
	}
	partitions := make([]partition, 0, goroutines)
	size := l / goroutines
	rem := l % goroutines
	low := 0
	high := size
	for i := 0; i < goroutines; i++ {
		partitions = append(partitions, partition{low, high})
		low = high
		high = high + size
		if i+rem+1 >= goroutines {
			high++
		}
	}
	return partitions
}
