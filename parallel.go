package stream

import "context"

type Parallel[E any, R any] struct {
	goroutines int
	slice      []E
	handler    func(index int, elem E) (isReturn bool, isComplete bool, result R)
}

func (p Parallel[E, R]) Run() []R {
	partitions := partition(p.slice, p.goroutines)
	resultChs := make([]chan []R, len(partitions))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, pa := range partitions {
		resultChs[i] = make(chan []R)
		go p.do(ctx, cancel, resultChs[i], pa)
	}

	result := p.resulted(resultChs, len(p.slice))
	return result
}

func (p Parallel[E, R]) do(
	ctx context.Context,
	cancel context.CancelFunc,
	resultCh chan []R,
	pa part) {

	defer close(resultCh)
	ret := make([]R, 0, pa.high-pa.low)

	for i := pa.low; i < pa.high; i++ {
		select {
		case <-ctx.Done():
			break
		default:
			isReturn, isComplete, r := p.handler(i, p.slice[i])
			if !isReturn {
				continue
			}
			ret = append(ret, r)
			if isComplete {
				cancel()
				break
			}
		}
	}

	if len(ret) > 0 {
		resultCh <- ret
	}
	return
}

func (p Parallel[E, R]) resulted(resultChs []chan []R, cap int) []R {
	results := make([]R, 0, cap)
	for _, resultCh := range resultChs {
		for result := range resultCh {
			results = append(results, result...)
		}
	}
	return results
}

// part  Uniform slices
// This selects a half-open range which includes the first element, but excludes the last.
type part struct {
	low  int //includes index
	high int //excludes index
}

// partition Given a specified source, evenly part according to the source.
func partition[E any](slice []E, goroutines int) []part {
	l := len(slice)
	if l == 0 {
		return nil
	}
	if goroutines > l {
		goroutines = l
	}
	partitions := make([]part, 0, goroutines)
	size := l / goroutines
	rem := l % goroutines
	low := 0
	high := size
	for i := 0; i < goroutines; i++ {
		partitions = append(partitions, part{low, high})
		low = high
		high = high + size
		if i+rem+1 >= goroutines {
			high++
		}
	}
	return partitions
}
