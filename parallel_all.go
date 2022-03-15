package stream

// ParallelAll All elements need to be processed in parallel, all return values are obtained, and then the parallel is ended.
// For SliceStream.Map, SliceStream.ForEach, SliceStream.Filter.
type ParallelAll[Elem any, TaskResult any] struct {
	// number of goroutine
	goroutines int
	slice      []Elem
	handler    ParallelHandler[Elem, TaskResult]
}

func (p ParallelAll[Elem, TaskResult]) Process(goroutines int, slice []Elem, handler ParallelHandler[Elem, TaskResult], _ ...TaskResult) []TaskResult {
	p.goroutines = goroutines
	p.slice = slice
	p.handler = handler

	partitions := partitionHandler(p.slice, p.goroutines)
	taskResultChs := make([]chan []TaskResult, len(partitions))
	for i, part := range partitions {
		taskResultChs[i] = make(chan []TaskResult)
		go p.task(taskResultChs[i], part)
	}

	result := p.resultHandler(taskResultChs)
	return result
}

func (p ParallelAll[Elem, TaskResult]) task(ch chan []TaskResult, part partition) {
	defer close(ch)
	ret := make([]TaskResult, 0, part.high-part.low)

	for i := part.low; i < part.high; i++ {
		isReturn, r := p.handler(i, p.slice[i])
		if isReturn {
			ret = append(ret, r)
		}
	}
	
	if len(ret) > 0 {
		ch <- ret
	}
	return
}

func (p ParallelAll[Elem, TaskResult]) resultHandler(taskResultChs []chan []TaskResult) []TaskResult {
	results := make([]TaskResult, 0, len(p.slice))
	for _, ch := range taskResultChs {
		for result := range ch {
			for _, r := range result {
				results = append(results, r)
			}
		}
	}
	return results
}
