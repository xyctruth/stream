package stream

import "context"

func parallel[Elem any, TaskResult any, Result any](cores int, slice []Elem, handler func(Elem) TaskResult, resultHandler func(results chan TaskResult) Result) Result {
	tasks := make(chan Elem, cores)
	results := make(chan TaskResult, cores)
	ctx, cancel := context.WithCancel(context.Background())

	parallelPublish(ctx, tasks, slice)
	parallelConsume[Elem, TaskResult](ctx, cores, tasks, results, handler)

	result := resultHandler(results)
	cancel()
	return result
}

func parallelPublish[Elem any](ctx context.Context, tasks chan Elem, slice []Elem) {
	go func() {
		for _, v := range slice {
			select {
			case <-ctx.Done():
			case tasks <- v:
			}
		}
	}()
}

func parallelConsume[Elem any, TaskResult any](ctx context.Context, cores int, tasks chan Elem, results chan TaskResult, handler func(Elem) TaskResult) {
	for i := 0; i < cores; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
				case task := <-tasks:
					results <- handler(task)
				}
			}
		}()
	}
}
