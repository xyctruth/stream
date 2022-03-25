package stream

func PipelineIntermediate[E any](slice []E, goroutines int, intermediate IntermediateStage[E]) []E {
	if slice == nil || intermediate == nil {
		return slice
	}

	if goroutines > 1 {
		handler := func(i int, v E) (isReturn bool, ret E) {
			return intermediate(i, v)
		}
		newSlice := ParallelProcess[E, E](
			goroutines,
			slice,
			handler,
			ParallelAllType)
		return newSlice
	}

	newSlice := make([]E, 0, len(slice))
	for i, v := range slice {
		isReturn, ret := intermediate(i, v)
		if !isReturn {
			continue
		}
		newSlice = append(newSlice, ret)
	}
	return newSlice
}

func PipelineTermination[E any, R any](slice []E, goroutines int, intermediate IntermediateStage[E], termination TerminationStage[E, R], parallelType ParallelType) []R {
	if goroutines > 1 {
		handler := func(i int, v E) (isReturn bool, ret R) {
			if intermediate != nil {
				isReturn, v = intermediate(i, v)
				if !isReturn {
					return
				}
			}
			isReturn, _, ret = termination(i, v)
			return
		}
		return ParallelProcess[E, R](
			goroutines,
			slice,
			handler,
			parallelType)
	}

	handler := func(i int, v E) (isReturn bool, isComplete bool, ret R) {
		if intermediate != nil {
			isReturn, v = intermediate(i, v)
			if !isReturn {
				return
			}
		}
		return termination(i, v)
	}

	newSlice := make([]R, 0, len(slice))
	for i, v := range slice {
		isReturn, isComplete, ret := handler(i, v)
		if !isReturn {
			continue
		}
		newSlice = append(newSlice, ret)
		if isComplete {
			return newSlice
		}
	}
	return newSlice
}

type IntermediateStage[E any] func(index int, e E) (bool, E)

func (s1 IntermediateStage[E]) Wrap(s2 IntermediateStage[E]) IntermediateStage[E] {
	if s1 == nil {
		return s2
	}
	return func(index int, e E) (isReturn bool, ret E) {
		isReturn, ret = s1(index, e)
		if !isReturn {
			return
		}
		return s2(index, ret)
	}
}

type TerminationStage[E any, R any] func(index int, e E) (isReturn bool, isComplete bool, ret R)
