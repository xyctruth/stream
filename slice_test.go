package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSliceStream(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "normal",
			input: []string{"a", "b"},
			want:  []string{"a", "b"},
		},
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).Slice()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceAt(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		index int
		want  int
	}{
		{
			name:  "normal",
			input: []int{1, 2, 3},
			index: 1,
			want:  2,
		},
		{
			name:  "normal2",
			input: []int{1, 2, 3},
			index: -1,
			want:  3,
		},
		{
			name:  "normal3",
			input: []int{1, 2, 3},
			index: 5,
			want:  0,
		},
		{
			name:  "normal4",
			input: []int{1, 2, 3},
			index: -4,
			want:  0,
		},
		{
			name:  "empty",
			input: []int{},
			index: 0,
			want:  0,
		},
		{
			name:  "nil",
			input: nil,
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).At(tt.index)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceAllMatch(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(v int) bool
		want      bool
	}{
		{
			name:      "all match",
			input:     []int{1, 2},
			predicate: func(v int) bool { return v < 3 },
			want:      true,
		},
		{
			name:      "no match",
			input:     []int{1, 2},
			predicate: func(v int) bool { return v > 3 },
			want:      false,
		},
		{
			name:      "empty",
			input:     []int{},
			predicate: func(v int) bool { return v > 3 },
			want:      true,
		},
		{
			name:      "nil",
			input:     nil,
			predicate: func(v int) bool { return v > 3 },
			want:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).AllMatch(tt.predicate)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceAnyMatch(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(v int) bool
		want      bool
	}{
		{
			name:      "match",
			input:     []int{1, 2},
			predicate: func(v int) bool { return v == 1 },
			want:      true,
		},
		{
			name:      "no match",
			input:     []int{1, 2},
			predicate: func(v int) bool { return v == 3 },
			want:      false,
		},
		{
			name:      "empty",
			input:     []int{},
			predicate: func(v int) bool { return v > 3 },
			want:      false,
		},
		{
			name:      "nil",
			input:     nil,
			predicate: func(v int) bool { return v > 3 },
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).AnyMatch(tt.predicate)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceCount(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "normal",
			input: []int{1, 2},
			want:  2,
		},
		{
			name:  "empty",
			input: []int{},
			want:  0,
		},
		{
			name:  "nil",
			input: nil,
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).Count()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceEqualFunc(t *testing.T) {
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
			got := NewSlice(tt.input).EqualFunc(tt.input2, func(a int, b int) bool { return a == b })
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		predicate func(v string) bool
		want      []string
	}{
		{
			name:      "match",
			input:     []string{"a", "b"},
			predicate: func(v string) bool { return v == "a" },
			want:      []string{"a"},
		},
		{
			name:      "no match",
			input:     []string{"a", "b"},
			predicate: func(v string) bool { return v == "c" },
			want:      []string{},
		},
		{
			name:      "nil",
			input:     nil,
			predicate: nil,
			want:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).Filter(tt.predicate)
			assert.Equal(t, tt.want, got.slice)
		})
	}
}

func TestSliceFirst(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "normal",
			input: []int{1, 2, 1},
			want:  1,
		},
		{
			name:  "empty",
			input: []int{},
			want:  0,
		},
		{
			name:  "nil",
			input: nil,
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).First()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceLimit(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		limit int
		want  []int
	}{
		{
			name:  "exceed",
			input: []int{1, 2, 1},
			limit: 5,
			want:  []int{1, 2, 1},
		},
		{
			name:  "no exceed",
			input: []int{1, 2, 1},
			limit: 2,
			want:  []int{1, 2},
		},
		{
			name:  "all",
			input: []int{1, 2, 1},
			limit: 3,
			want:  []int{1, 2, 1},
		},
		{
			name:  "limit(0)",
			input: []int{1, 2, 1},
			limit: 0,
			want:  []int{},
		},
		{
			name:  "empty",
			input: []int{},
			limit: 5,
			want:  []int{},
		},
		{
			name:  "nil",
			input: nil,
			limit: 5,
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).Limit(tt.limit)
			assert.Equal(t, tt.want, got.slice)
		})
	}
}

func TestSliceMap(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		mapper func(int) int
		want   []int
	}{
		{
			name:   "normal",
			input:  []int{1, 2, 1},
			mapper: func(i int) int { return i * 2 },
			want:   []int{2, 4, 2},
		},
		{
			name:   "empty",
			input:  []int{},
			mapper: func(i int) int { return i * 2 },
			want:   []int{},
		},
		{
			name:   "nil",
			input:  nil,
			mapper: func(i int) int { return i * 2 },
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).Map(tt.mapper)
			assert.Equal(t, tt.want, got.slice)
		})
	}
}

func TestSliceReduce(t *testing.T) {
	tests := []struct {
		name        string
		input       []int
		accumulator func(int, int) int
		want        int
	}{
		{
			name:        "normal",
			input:       []int{1, 2, 1, 10},
			accumulator: func(i int, j int) int { return i + j },
			want:        14,
		},
		{
			name:        "empty",
			input:       []int{},
			accumulator: func(i int, j int) int { return i + j },
			want:        0,
		},
		{
			name:        "nil",
			input:       nil,
			accumulator: func(i int, j int) int { return i + j },
			want:        0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).Reduce(tt.accumulator)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceSortFunc(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		less  func(a, b int) bool
		want  []int
	}{
		{
			name:  "normal",
			input: []int{1, 2, 1, 5},
			less:  func(a, b int) bool { return a > b },
			want:  []int{5, 2, 1, 1},
		},
		{
			name:  "normal",
			input: []int{1, 2, 1, 5},
			less:  func(a, b int) bool { return a < b },
			want:  []int{1, 1, 2, 5},
		},
		{
			name:  "empty",
			input: []int{},
			less:  func(a, b int) bool { return a > b },
			want:  []int{},
		},
		{
			name:  "nil",
			input: nil,
			less:  func(a, b int) bool { return a > b },
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).SortFunc(tt.less)
			assert.Equal(t, tt.want, got.slice)
		})
	}
}

func TestSliceSortStableFunc(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		less  func(a, b int) bool
		want  []int
	}{
		{
			name:  "normal",
			input: []int{1, 2, 1, 5},
			less:  func(a, b int) bool { return a > b },
			want:  []int{5, 2, 1, 1},
		},
		{
			name:  "normal",
			input: []int{1, 2, 1, 5},
			less:  func(a, b int) bool { return a < b },
			want:  []int{1, 1, 2, 5},
		},
		{
			name:  "empty",
			input: []int{},
			less:  func(a, b int) bool { return a > b },
			want:  []int{},
		},
		{
			name:  "nil",
			input: nil,
			less:  func(a, b int) bool { return a > b },
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).SortStableFunc(tt.less)
			assert.Equal(t, tt.want, got.slice)
		})
	}
}
