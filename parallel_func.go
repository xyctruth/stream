package stream

// singleResultHandler A handler for parallelProcess of a single return value.
// Get the first return value, and the parallel processing ends.
// For sliceStream.AllMatch, sliceStream.AnyMatch, sliceStream.FindFunc.
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

// multipleResultHandler A handler for parallelProcess of a multiple return value.
// Need to get the results of all parallel processing.
// For sliceStream.Map, sliceStream.ForEach, sliceStream.Filter.
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

// partitionHandler Given a specified slice, evenly partition according to the slice.
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
