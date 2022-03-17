package stream

// ParallelAll All elements need to be processed in parallel, all return values are obtained, and then the parallel is ended.
//
// - E elements type
// - R result type
//
// For SliceStream.Map, SliceStream.Filter.
type ParallelAll[E any, R any] struct {
	slice   []E                   // element to be processed
	handler ParallelHandler[E, R] // handler function
}

func (p ParallelAll[E, R]) Process(goroutines int, slice []E, handler ParallelHandler[E, R]) []R {
	p.slice = slice
	p.handler = handler

	partitions := partition(p.slice, goroutines)
	resultChs := make([]chan []R, len(partitions))

	for i, pa := range partitions {
		resultChs[i] = make(chan []R)
		go p.do(resultChs[i], pa)
	}

	result := p.resulted(resultChs)
	return result
}

func (p ParallelAll[_, R]) do(resultCh chan []R, pa part) {
	defer close(resultCh)
	ret := make([]R, 0, pa.high-pa.low)

	for i := pa.low; i < pa.high; i++ {
		isReturn, r := p.handler(i, p.slice[i])
		if isReturn {
			ret = append(ret, r)
		}
	}

	if len(ret) > 0 {
		resultCh <- ret
	}
	return
}

func (p ParallelAll[_, R]) resulted(resultChs []chan []R) []R {
	results := make([]R, 0, len(p.slice))
	for _, resultCh := range resultChs {
		for result := range resultCh {
			for _, r := range result {
				results = append(results, r)
			}
		}
	}
	return results
}
