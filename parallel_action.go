package stream

import "sync"

// ParallelAction All elements need to be processed in parallel, no return value required
//
// - E elements type
// - R result type
//
// For  SliceStream.ForEach
type ParallelAction[E any, R any] struct {
}

func (p ParallelAction[E, R]) Process(goroutines int, slice []E, handler ParallelHandler[E, R]) []R {
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
