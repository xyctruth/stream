package stream

// partition  Uniform slices
// This selects a half-open range which includes the first element, but excludes the last.
type partition struct {
	low  int //includes index
	high int //excludes index
}

// partitionHandler Given a specified slice, evenly partition according to the slice.
func partitionHandler[Elem any](slice []Elem, goroutines int) []partition {
	l := len(slice)
	if l == 0 {
		return nil
	}
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
