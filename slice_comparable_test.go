package stream

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceComparableDistinct(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "has distinct",
			input: []int{1, 2, 1},
			want:  []int{1, 2},
		},
		{
			name:  "no distinct",
			input: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
		{
			name:  "empty",
			input: []int{},
			want:  []int{},
		},
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSliceByComparable(tt.input).Distinct().ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).Distinct().ToSlice()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceComparableEqual(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		input2 []int
		want   bool
	}{
		{
			name:   "normal",
			input:  []int{1, 2},
			input2: []int{1, 2},
			want:   true,
		},
		{
			name:   "normal",
			input:  []int{1, 2},
			input2: []int{1, 2},
			want:   true,
		},
		{
			name:   "empty",
			input:  []int{},
			input2: []int{},
			want:   true,
		},
		{
			name:   "nil",
			input:  nil,
			input2: nil,
			want:   true,
		},
		{
			name:   "nil and empty",
			input:  []int{},
			input2: nil,
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSliceByComparable(tt.input).Equal(tt.input2)
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).Equal(tt.input2)
			assert.Equal(t, tt.want, got)
		})
	}
}
