package stream

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipelineStages(t *testing.T) {
	p := &Pipeline[int]{
		source: []int{1, 2, 3},
	}

	p.AddStage(func(index int, e int) (isReturn bool, isComplete bool, ret int) {
		return true, false, e * 10
	})

	isReturn, isComplete, ret := p.stages(0, 1)
	assert.Equal(t, true, isReturn)
	assert.Equal(t, false, isComplete)
	assert.Equal(t, ret, 10)

	p.AddStage(func(index int, e int) (isReturn bool, isComplete bool, ret int) {
		return true, false, e + 10
	})
	isReturn, isComplete, ret = p.stages(0, 1)
	assert.Equal(t, true, isReturn)
	assert.Equal(t, false, isComplete)
	assert.Equal(t, ret, 20)

	p.evaluation()
	assert.Equal(t, p.source, []int{20, 30, 40})
	assert.Nil(t, p.stages)
}

func TestPipeByTermination(t *testing.T) {
	tests := []struct {
		name       string
		goroutines int
	}{
		{
			name:       "case",
			goroutines: 0,
		},
		{
			name:       "case",
			goroutines: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pipeline[int]{
				source:     []int{1, 2, 3},
				goroutines: tt.goroutines,
			}
			p.AddStage(func(index int, e int) (isReturn bool, isComplete bool, ret int) {
				return true, false, e * 10
			})

			rets := pipelineTermination(p, func(index int, e int) (isReturn bool, isComplete bool, ret int) {
				if index == 1 {
					return true, true, e * 10
				}
				return false, false, e * 10
			})

			assert.Equal(t, []int{200}, rets)
			assert.Nil(t, p.stages)
		})
	}
}
