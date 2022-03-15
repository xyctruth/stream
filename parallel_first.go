package stream

import (
	"context"
	"sync"
)

// ParallelFirst parallel processing ends as soon as the first return value is obtained.
// For SliceStream.AllMatch, SliceStream.AnyMatch, SliceStream.FindFunc.
type ParallelFirst[Elem any, TaskResult any] struct {
	// number of goroutine
	goroutines int
	slice      []Elem
	handler    ParallelHandler[Elem, TaskResult]
}

func (p ParallelFirst[Elem, TaskResult]) Process(goroutines int, slice []Elem, handler ParallelHandler[Elem, TaskResult], defaultResults ...TaskResult) []TaskResult {
	p.goroutines = goroutines
	p.slice = slice
	p.handler = handler

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	partitions := partitionHandler(p.slice, p.goroutines)

	wg := sync.WaitGroup{}
	wg.Add(len(partitions))
	taskResultCh := make(chan TaskResult, len(partitions))

	for _, part := range partitions {
		go p.task(ctx, &wg, taskResultCh, part)
	}
	go func() {
		wg.Wait()
		close(taskResultCh)
	}()
	result := p.resultHandler(taskResultCh, defaultResults...)
	return result
}

func (p ParallelFirst[Elem, TaskResult]) task(ctx context.Context, wg *sync.WaitGroup, taskResultCh chan TaskResult, part partition) {
	defer wg.Done()
	for i := part.low; i < part.high; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			isReturn, r := p.handler(i+part.low, p.slice[i])
			// Get the first matching value, end the parallel
			if isReturn {
				taskResultCh <- r
				return
			}
		}
	}
}

func (p ParallelFirst[Elem, TaskResult]) resultHandler(taskResultCh chan TaskResult, defaultResults ...TaskResult) []TaskResult {
	for result := range taskResultCh {
		return []TaskResult{result}
	}
	return defaultResults
}
