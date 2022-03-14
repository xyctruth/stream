package stream

import (
	"context"
)

type parallel[Elem any, TaskResult any, Result any] struct {
	// number of goroutine
	goroutines int
	slice      []Elem
	handler    parallelHandler[Elem, TaskResult]

	// Pass parallel processing results
	taskResultChs []chan []TaskResult
	// Process the  parallel returned result
	// implement singleResultHandler, multipleResultHandler
	resultHandler parallelResultHandler[TaskResult, Result]

	// Whether to wait for all elements to complete
	// true:  All elements need to be processed. Collaborate singleResultHandler
	// false: Get the first matching value, end the parallel. Collaborate multipleResultHandler
	isWaitAllDone bool
}

type parallelHandler[Elem any, TaskResult any] func(index int, elem Elem) (isReturn bool, taskResult TaskResult)

type parallelResultHandler[TaskResult any, Result any] func(taskResultChs []chan []TaskResult) Result

// parallelProcess Provides parallel processing capabilities.
func parallelProcess[Elem any, TaskResult any, Result any](
	goroutines int,
	slice []Elem,
	handler parallelHandler[Elem, TaskResult],
	resultHandler parallelResultHandler[TaskResult, Result],
	isWaitAllDone bool) Result {

	p := parallel[Elem, TaskResult, Result]{
		goroutines:    goroutines,
		slice:         slice,
		handler:       handler,
		resultHandler: resultHandler,
		isWaitAllDone: isWaitAllDone,
	}
	return p.process()
}

func (p parallel[Elem, TaskResult, Result]) process() Result {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(p.slice) > 0 {
		partitions := partitionHandler(p.slice, p.goroutines)
		p.taskResultChs = make([]chan []TaskResult, len(partitions))
		for i, part := range partitions {
			p.taskResultChs[i] = make(chan []TaskResult)
			go p.task(ctx, cancel, p.taskResultChs[i], part)
		}
	}

	result := p.resultHandler(p.taskResultChs)
	return result
}

func (p parallel[Elem, TaskResult, Result]) task(ctx context.Context, cancel context.CancelFunc, ch chan []TaskResult, part partition) {
	defer close(ch)

	//  All elements need to be processed
	if p.isWaitAllDone {
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

	// Get the first matching value, end the parallel
	for i := part.low; i < part.high; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			isReturn, r := p.handler(i+part.low, p.slice[i])
			if isReturn {
				ch <- []TaskResult{r}
				cancel()
				return
			}
		}
	}
}
