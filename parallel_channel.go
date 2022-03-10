package stream

import "context"

type parallelChannel[Elem any, TaskResult any, Result any] struct {
	goroutines    int
	slice         []Elem
	handler       func(Elem) TaskResult
	resultHandler func(results chan TaskResult) Result
	tasks         chan Elem
	results       chan TaskResult
}

func (p parallelChannel[Elem, TaskResult, Result]) process() Result {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p.publish(ctx)
	p.consume(ctx)
	result := p.resultHandler(p.results)
	return result
}

func (p parallelChannel[Elem, TaskResult, Result]) publish(ctx context.Context) {
	go func() {
		for _, v := range p.slice {
			select {
			case <-ctx.Done():
				return
			case p.tasks <- v:
			}
		}
	}()
}

func (p parallelChannel[Elem, TaskResult, Result]) consume(ctx context.Context) {
	for i := 0; i < p.goroutines; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case task := <-p.tasks:
					p.results <- p.handler(task)
				}
			}
		}()
	}
}
