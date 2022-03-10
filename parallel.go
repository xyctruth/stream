package stream

import "context"

type taskResult[V any] struct {
	v         V
	predicate bool
}

type parallelHandler[Elem any, TaskResult any] func(int, Elem) taskResult[TaskResult]

type parallelResultHandler[TaskResult any, Result any] func(results chan []taskResult[TaskResult]) Result

type parallel[Elem any, TaskResult any, Result any] struct {
	goroutines    int
	slice         []Elem
	handler       parallelHandler[Elem, TaskResult]
	resultHandler parallelResultHandler[TaskResult, Result]
	taskResultCh  chan []taskResult[TaskResult]
	isWaitAllDone bool
}

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
		taskResultCh:  make(chan []taskResult[TaskResult], goroutines),
		isWaitAllDone: isWaitAllDone,
	}
	return p.process()
}

func (p parallel[Elem, TaskResult, Result]) process() Result {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(p.slice) > 0 {
		partition, size := p.partition(p.slice, p.goroutines)
		for i, s := range partition {
			go p.task(ctx, s, i*size)
		}
	}

	result := p.resultHandler(p.taskResultCh)
	return result
}

func (p parallel[Elem, TaskResult, Result]) task(ctx context.Context, slice []Elem, offset int) {
	if p.isWaitAllDone {
		ret := make([]taskResult[TaskResult], 0, len(slice))
		for i, elem := range slice {
			ret = append(ret, p.handler(i+offset, elem))
		}
		p.taskResultCh <- ret
		return
	}

	for i, elem := range slice {
		select {
		case <-ctx.Done():
			return
		default:
			p.taskResultCh <- []taskResult[TaskResult]{p.handler(i, elem)}
		}
	}
}

func (p parallel[Elem, TaskResult, Result]) partition(slice []Elem, cores int) ([][]Elem, int) {
	var ret [][]Elem
	l := len(slice)

	if cores > l {
		cores = l
	}

	size := int(float64(l) / float64(cores))
	rem := l % cores
	for i := 0; i < cores; i++ {
		s := i * size
		e := (i + 1) * size
		if i == cores-1 {
			e = e + rem
		}
		ret = append(ret, slice[s:e])
	}
	return ret, size
}

func parallelResultHandlerMatch(c bool, count int) parallelResultHandler[bool, bool] {
	return func(taskResultCh chan []taskResult[bool]) bool {
		for i := 0; i < count; {
			result := <-taskResultCh
			for _, r := range result {
				if r.predicate == c {
					return c
				}
			}
			i = i + len(result)
		}
		return !c
	}
}

func parallelResultHandlerEach[Elem any](count int) parallelResultHandler[Elem, []Elem] {
	return func(results chan []taskResult[Elem]) []Elem {
		newSlice := make([]Elem, 0, count)
		for i := 0; i < count; {
			result := <-results
			for _, r := range result {
				if r.predicate {
					newSlice = append(newSlice, r.v)
				}
			}
			i = i + len(result)
		}
		return newSlice
	}
}
