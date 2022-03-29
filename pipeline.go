package stream

type Stage[E any, R any] func(index int, e E) (isReturn bool, isComplete bool, ret R)

type Pipeline[E any] struct {
	source     []E
	goroutines int
	stages     Stage[E, E]
}

func (pipe *Pipeline[E]) AddStage(s2 Stage[E, E]) {
	if pipe.stages == nil {
		pipe.stages = s2
		return
	}
	s1 := pipe.stages
	pipe.stages = func(index int, e E) (isReturn bool, isComplete bool, ret E) {
		isReturn, _, ret = s1(index, e)
		if !isReturn {
			return
		}
		return s2(index, ret)
	}
}

func (pipe *Pipeline[E]) evaluation() {
	if pipe.source == nil || pipe.stages == nil {
		return
	}
	defer func() {
		pipe.stages = nil
	}()
	pipe.source = pipelineRun(pipe.source, pipe.goroutines, pipe.stages)
}

func pipelineTermination[E any, R any](pipe *Pipeline[E], terminationStage Stage[E, R]) []R {
	if pipe.source == nil {
		return nil
	}

	defer func() {
		pipe.stages = nil
	}()

	var stages Stage[E, R]
	if pipe.stages == nil {
		stages = terminationStage
	} else {
		stages = func(i int, v E) (isReturn bool, isComplete bool, ret R) {
			isReturn, _, v = pipe.stages(i, v)
			if !isReturn {
				return
			}
			isReturn, isComplete, ret = terminationStage(i, v)
			return
		}
	}
	return pipelineRun(pipe.source, pipe.goroutines, stages)
}

func pipelineRun[E any, R any](source []E, goroutines int, stages Stage[E, R]) []R {
	if goroutines > 1 {
		return Parallel[E, R]{goroutines, source, stages}.Run()
	}

	newSlice := make([]R, 0, len(source))
	for i, v := range source {
		isReturn, isComplete, ret := stages(i, v)
		if !isReturn {
			continue
		}
		newSlice = append(newSlice, ret)
		if isComplete {
			break
		}
	}
	return newSlice
}
