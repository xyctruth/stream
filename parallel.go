package stream

import (
	"context"
	"sync"
)

type parallel[Elem any, TaskResult any, Result any] struct {
	goroutines    int
	slice         []Elem
	handler       func(int, Elem) (isReturn bool, taskResult TaskResult)
	resultHandler parallelResultHandler[TaskResult, Result]
	taskResultCh  chan []TaskResult
	isWaitAllDone bool
}

type parallelHandler[Elem any, TaskResult any] func(int, Elem) (isReturn bool, taskResult TaskResult)

type parallelResultHandler[TaskResult any, Result any] func(taskResultCh chan []TaskResult) Result

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
		taskResultCh:  make(chan []TaskResult, goroutines),
		isWaitAllDone: isWaitAllDone,
	}
	return p.process()
}

func (p parallel[Elem, TaskResult, Result]) process() Result {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	if len(p.slice) > 0 {
		partitions, size := partition(p.slice, p.goroutines)

		wg.Add(len(partitions))
		for i, s := range partitions {
			go p.task(ctx, &wg, s, i*size)
		}
	}

	go func() {
		wg.Wait()
		close(p.taskResultCh)
	}()

	result := p.resultHandler(p.taskResultCh)
	return result
}

func (p parallel[Elem, TaskResult, Result]) task(ctx context.Context, wg *sync.WaitGroup, slice []Elem, offset int) {
	defer wg.Done()

	if p.isWaitAllDone {
		ret := make([]TaskResult, 0, len(slice))
		for i, elem := range slice {
			isReturn, r := p.handler(i+offset, elem)
			if isReturn {
				ret = append(ret, r)
			}
		}
		if len(ret) > 0 {
			p.taskResultCh <- ret
		}
		return
	}

	for i, elem := range slice {
		select {
		case <-ctx.Done():
			return
		default:
			isReturn, r := p.handler(i+offset, elem)
			if isReturn {
				p.taskResultCh <- []TaskResult{r}
			}
		}
	}
}
