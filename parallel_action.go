package stream

import "sync"

// ParallelAction All elements need to be processed in parallel, no return value required
//
// For  SliceStream.ForEach
type ParallelAction[Elem any, Result any] struct {
}

func (p ParallelAction[Elem, Result]) Process(goroutines int, slice []Elem, handler ParallelHandler[Elem, Result]) []Result {
	wg := sync.WaitGroup{}
	partitions := partition(slice, goroutines)
	wg.Add(len(partitions))
	for _, pa := range partitions {
		go func(pa part) {
			defer wg.Done()
			for i := pa.low; i < pa.high; i++ {
				handler(i, slice[i])
			}
		}(pa)
	}
	wg.Wait()
	return nil
}
