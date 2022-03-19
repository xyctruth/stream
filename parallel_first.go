package stream

import (
	"context"
	"sync"
)

// ParallelFirst parallel processing ends as soon as the first return value is obtained.
//
// - E elements type
// - R result type
//
// For SliceStream.AllMatch, SliceStream.AnyMatch, SliceStream.FindFunc.
type ParallelFirst[E any, R any] struct {
	slice   []E                       // element to be processed
	handler ParallelHandlerFunc[E, R] // handler function
}

func (p ParallelFirst[E, R]) Process(goroutines int, slice []E, handler ParallelHandlerFunc[E, R]) []R {
	p.slice = slice
	p.handler = handler

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	partitions := partition(p.slice, goroutines)
	resultCh := make(chan R, len(partitions))
	wg.Add(len(partitions))

	for _, pa := range partitions {
		go p.do(ctx, &wg, resultCh, pa)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	result := p.resulted(resultCh)
	return result
}

func (p ParallelFirst[_, R]) do(ctx context.Context, wg *sync.WaitGroup, resultCh chan R, pa part) {
	defer wg.Done()
	for i := pa.low; i < pa.high; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			isReturn, r := p.handler(i+pa.low, p.slice[i])
			// Get the first matching value, end the parallel
			if isReturn {
				resultCh <- r
				return
			}
		}
	}
}

func (p ParallelFirst[_, R]) resulted(resultCh chan R) []R {
	for result := range resultCh {
		return []R{result}
	}
	return nil
}
