package stream

import (
	"context"
	"sync"
)

// ParallelFirst parallel processing ends as soon as the first return value is obtained.
//
// For SliceStream.AllMatch, SliceStream.AnyMatch, SliceStream.FindFunc.
type ParallelFirst[Elem any, Result any] struct {
	slice   []Elem                        // element to be processed
	handler ParallelHandler[Elem, Result] // handler function
}

func (p ParallelFirst[Elem, Result]) Process(goroutines int, slice []Elem, handler ParallelHandler[Elem, Result]) []Result {
	p.slice = slice
	p.handler = handler

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	partitions := partition(p.slice, goroutines)
	resultCh := make(chan Result, len(partitions))
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

func (p ParallelFirst[_, Result]) do(ctx context.Context, wg *sync.WaitGroup, taskResultCh chan Result, pa part) {
	defer wg.Done()
	for i := pa.low; i < pa.high; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			isReturn, r := p.handler(i+pa.low, p.slice[i])
			// Get the first matching value, end the parallel
			if isReturn {
				taskResultCh <- r
				return
			}
		}
	}
}

func (p ParallelFirst[_, TaskResult]) resulted(resultCh chan TaskResult) []TaskResult {
	for result := range resultCh {
		return []TaskResult{result}
	}
	return nil
}
