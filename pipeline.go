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
	pipe.source = pipelineRun(pipe, pipe.stages)
}

func (pipe *Pipeline[E]) evaluationBool(terminal Stage[E, bool]) (ret []bool) {
	ret = pipelineRun(pipe, wrapTerminal[E, bool](pipe.stages, terminal))
	return ret
}

func (pipe *Pipeline[E]) evaluationInt(terminal Stage[E, int]) (ret []int) {
	ret = pipelineRun(pipe, wrapTerminal[E, int](pipe.stages, terminal))
	return ret
}

func wrapTerminal[E any, R any](stage Stage[E, E], terminalStage Stage[E, R]) Stage[E, R] {
	var stages Stage[E, R]
	if stage == nil {
		stages = terminalStage
	} else {
		stages = func(i int, v E) (isReturn bool, isComplete bool, ret R) {
			isReturn, _, v = stage(i, v)
			if !isReturn {
				return
			}
			isReturn, isComplete, ret = terminalStage(i, v)
			return
		}
	}
	return stages
}

func pipelineRun[E any, R any](pipe *Pipeline[E], stages Stage[E, R]) []R {
	defer func() {
		pipe.stages = nil
	}()

	if pipe.goroutines > 1 {
		return Parallel[E, R]{pipe.goroutines, pipe.source, stages}.Run()
	}

	results := make([]R, 0, len(pipe.source))
	for i, v := range pipe.source {
		isReturn, isComplete, ret := stages(i, v)
		if !isReturn {
			continue
		}
		results = append(results, ret)
		if isComplete {
			break
		}
	}
	return results
}
