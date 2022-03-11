package stream

func parallelResultHandlerMatch(defaultVal bool) parallelResultHandler[bool, bool] {
	return func(taskResultCh chan []bool) bool {
		for result := range taskResultCh {
			for _, r := range result {
				return r
			}
		}
		return defaultVal
	}
}

func parallelResultHandlerEach[Elem any](count int) parallelResultHandler[Elem, []Elem] {
	return func(taskResultCh chan []Elem) []Elem {
		newSlice := make([]Elem, 0, count)
		for result := range taskResultCh {
			for _, r := range result {
				newSlice = append(newSlice, r)
			}
		}
		return newSlice
	}
}

func partition[Elem any](slice []Elem, goroutines int) ([][]Elem, int) {
	var ret [][]Elem
	l := len(slice)

	if goroutines > l {
		goroutines = l
	}

	size := int(float64(l) / float64(goroutines))
	rem := l % goroutines
	for i := 0; i < goroutines; i++ {
		s := i * size
		e := (i + 1) * size
		if i == goroutines-1 {
			e = e + rem
		}
		ret = append(ret, slice[s:e])
	}
	return ret, size
}
