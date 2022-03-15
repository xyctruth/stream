package stream

// ParallelAll All elements need to be processed in parallel, all return values are obtained, and then the parallel is ended.
//
// For SliceStream.Map, SliceStream.Filter.
type ParallelAll[Elem any, Result any] struct {
	slice   []Elem                        // element to be processed
	handler ParallelHandler[Elem, Result] // handler function
}

func (p ParallelAll[Elem, Result]) Process(goroutines int, slice []Elem, handler ParallelHandler[Elem, Result]) []Result {
	p.slice = slice
	p.handler = handler

	partitions := partition(p.slice, goroutines)
	resultChs := make([]chan []Result, len(partitions))

	for i, pa := range partitions {
		resultChs[i] = make(chan []Result)
		go p.do(resultChs[i], pa)
	}

	result := p.resulted(resultChs)
	return result
}

func (p ParallelAll[_, Result]) do(resultCh chan []Result, pa part) {
	defer close(resultCh)
	ret := make([]Result, 0, pa.high-pa.low)

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

func (p ParallelAll[_, Result]) resulted(resultChs []chan []Result) []Result {
	results := make([]Result, 0, len(p.slice))
	for _, resultCh := range resultChs {
		for result := range resultCh {
			for _, r := range result {
				results = append(results, r)
			}
		}
	}
	return results
}
