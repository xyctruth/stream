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

type partition[Elem any] struct {
	slice      []Elem
	startIndex int
}

func partitionHandler[Elem any](slice []Elem, goroutines int) []partition[Elem] {
	l := len(slice)
	if goroutines > l {
		goroutines = l
	}
	partitions := make([]partition[Elem], 0, goroutines)

	size := int(float64(l) / float64(goroutines))
	rem := l % goroutines
	for i := 0; i < goroutines; i++ {
		s := i * size
		e := (i + 1) * size
		if i == goroutines-1 {
			e = e + rem
		}
		partitions = append(partitions, partition[Elem]{slice: slice[s:e], startIndex: size * i})
	}
	return partitions
}
