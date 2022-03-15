package stream

type ParallelHandler[Elem any, TaskResult any] func(index int, elem Elem) (isReturn bool, taskResult TaskResult)

// Parallel Provides parallel processing capabilities.
//
// ParallelFirst parallel processing ends as soon as the first return value is obtained.
// ParallelAll All elements need to be processed in parallel, all return values are obtained, and then the parallel is ended.
type Parallel[Elem any, TaskResult any] interface {
	Process(goroutines int, slice []Elem, handler ParallelHandler[Elem, TaskResult], defaultResult ...TaskResult) []TaskResult
}

func ParallelProcess[T Parallel[Elem, TaskResult], Elem any, TaskResult any](
	goroutines int,
	slice []Elem,
	handler ParallelHandler[Elem, TaskResult],
	defaultResults ...TaskResult) []TaskResult {

	var p T
	return p.Process(goroutines, slice, handler, defaultResults...)
}
