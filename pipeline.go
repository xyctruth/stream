package stream

type Stage[E any, R any] func(index int, e E) (isReturn bool, isComplete bool, ret R)

type Pipe[E any] struct {
	source     []E
	goroutines int
	stages     Stage[E, E]
}

func (pipe *Pipe[E]) AddStage(s2 Stage[E, E]) {
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

func (pipe *Pipe[E]) Run() {
	if pipe.source == nil || pipe.stages == nil {
		return
	}
	defer func() {
		pipe.stages = nil
	}()

	if pipe.goroutines > 1 {
		pipe.source = Parallel[E, E]{
			pipe.goroutines,
			pipe.source,
			pipe.stages}.Run()
		return
	}

	newSlice := make([]E, 0, len(pipe.source))
	for i, v := range pipe.source {
		isReturn, _, ret := pipe.stages(i, v)
		if !isReturn {
			continue
		}
		newSlice = append(newSlice, ret)
	}
	pipe.source = newSlice
}

func PipeByTermination[E any, R any](pipe *Pipe[E], terminationStage Stage[E, R]) []R {
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

	if pipe.goroutines > 1 {
		return Parallel[E, R]{
			pipe.goroutines,
			pipe.source,
			stages}.Run()
	}

	newSlice := make([]R, 0, len(pipe.source))
	for i, v := range pipe.source {
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
