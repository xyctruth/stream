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

	size := l / goroutines
	rem := l % goroutines
	var rg = [2]int{0, size}
	for i := 0; i < goroutines; i++ {
		partitions = append(partitions, partition[Elem]{slice: slice[rg[0]:rg[1]], startIndex: rg[0]})
		rg[0] = rg[1];
		rg[1] = rg[0] + size;
		if i + rem + 1 >= goroutines{
			rg[1] = rg[1] + 1;
		}
	}
	return partitions
}
