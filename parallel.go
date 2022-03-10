package stream

import "context"

type parallel[Elem any, TaskResult any, Result any] struct {
	goroutines    int
	slice         []Elem
	handler       func(Elem) TaskResult
	resultHandler func(results chan []TaskResult) Result
	tasks         chan Elem
	results       chan []TaskResult
	isWaitAllDone bool
}

func parallelProcess[Elem any, TaskResult any, Result any](goroutines int, slice []Elem, handler func(Elem) TaskResult, resultHandler func(results chan []TaskResult) Result) Result {
	p := parallel[Elem, TaskResult, Result]{
		goroutines:    goroutines,
		slice:         slice,
		handler:       handler,
		resultHandler: resultHandler,
		tasks:         make(chan Elem, goroutines),
		results:       make(chan []TaskResult, goroutines),
	}
	return p.process()
}

func (p parallel[Elem, TaskResult, Result]) process() Result {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(p.slice) > 0 {
		partition := p.partition(p.slice, p.goroutines)

		for _, s := range partition {
			go p.consume(ctx, s)
		}
	}
	result := p.resultHandler(p.results)
	return result
}

func (p parallel[Elem, TaskResult, Result]) consume(ctx context.Context, slice []Elem) {
	if p.isWaitAllDone {
		var ret []TaskResult
		for _, elem := range slice {
			ret = append(ret, p.handler(elem))
		}
		p.results <- ret
		return
	}

	for _, elem := range slice {
		select {
		case <-ctx.Done():
			return
		default:
			p.results <- []TaskResult{p.handler(elem)}
		}
	}
}

func (p parallel[Elem, TaskResult, Result]) partition(slice []Elem, cores int) [][]Elem {
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
	return ret
}
