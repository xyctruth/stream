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
			name:  "case",
			input: []int{1, 2, 1},
			want:  []int{1, 2},
		},
		{
			name:  "case",
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
			name:   "case",
			input:  []int{1, 2},
			input2: []int{1, 2},
			want:   true,
		},
		{
			name:   "case",
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

func TestSliceComparableFind(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		dest  int
		want  int
	}{
		{
			name:  "case",
			input: []int{1, 2, 1, 2, 1},
			dest:  1,
			want:  0,
		},
		{
			name:  "case",
			input: []int{1, 2, 1, 2, 1},
			dest:  2,
			want:  1,
		},
		{
			name:  "case",
			input: []int{1, 2},
			dest:  3,
			want:  -1,
		},
		{
			name:  "empty",
			input: []int{},
			dest:  1,
			want:  -1,
		},
		{
			name:  "nil",
			input: nil,
			dest:  1,
			want:  -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSliceByComparable(tt.input).Find(tt.dest)
			assert.Equal(t, tt.want, got)
		})
	}
}
