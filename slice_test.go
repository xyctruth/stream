package stream

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func newArray(count int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	s := make([]int, count)
	for i := 0; i < count; i++ {
		s[i] = r.Intn(count * 2)
	}
	return s
}

func newArrayN(count int, n int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	s := make([]int, count)
	for i := 0; i < count; i++ {
		s[i] = r.Intn(n)
	}
	return s
}

func TestNewSliceStream(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "case",
			input: []string{"a", "b"},
			want:  []string{"a", "b"},
		},
		{
			name:  "empty",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).ToSlice()
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
			name:  "case",
			input: []int{1, 2, 3},
			index: 1,
			want:  2,
		},
		{
			name:  "case",
			input: []int{1, 2, 3},
			index: -1,
			want:  3,
		},
		{
			name:  "case",
			input: []int{1, 2, 3},
			index: 5,
			want:  0,
		},
		{
			name:  "case",
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
			name:      "case",
			input:     []int{1, 2},
			predicate: func(v int) bool { return v < 3 },
			want:      true,
		},
		{
			name:      "case",
			input:     []int{1, 2},
			predicate: func(v int) bool { return v > 3 },
			want:      false,
		},
		{
			name:      "case",
			input:     newArrayN(100, 200),
			predicate: func(v int) bool { return v > 100 },
			want:      false,
		},
		{
			name:      "case",
			input:     newArrayN(100, 200),
			predicate: func(v int) bool { return v < 200 },
			want:      true,
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

			got = NewSlice(tt.input).Parallel(2).AllMatch(tt.predicate)
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
			name:      "case",
			input:     []int{1, 2},
			predicate: func(v int) bool { return v == 1 },
			want:      true,
		},
		{
			name:      "case",
			input:     []int{1, 2},
			predicate: func(v int) bool { return v == 3 },
			want:      false,
		},
		{
			name:      "case",
			input:     append([]int{300}, newArrayN(100, 200)...),
			predicate: func(v int) bool { time.Sleep(time.Millisecond); return v == 300 },
			want:      true,
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

			got = NewSlice(tt.input).Parallel(2).AnyMatch(tt.predicate)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceAppend(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		input2 []int
		want   []int
	}{
		{
			name:   "case",
			input:  []int{1, 2},
			input2: []int{3, 4},
			want:   []int{1, 2, 3, 4},
		},
		{
			name:   "empty",
			input:  []int{1, 2},
			input2: []int{},
			want:   []int{1, 2},
		},
		{
			name:   "nil",
			input:  []int{1, 2},
			input2: nil,
			want:   []int{1, 2},
		},
		{
			name:   "empty",
			input:  []int{},
			input2: []int{3, 4},
			want:   []int{3, 4},
		},
		{
			name:   "nil",
			input:  nil,
			input2: []int{3, 4},
			want:   []int{3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).Append(tt.input2...).ToSlice()
			assert.Equal(t, tt.want, got)

			got[0] = 100000
			if len(tt.input) > 0 && len(got) > 0 {
				assert.NotEqual(t, tt.want[0], got[0])
			}

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
			name:  "case",
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
			got := NewSlice(tt.input).EqualFunc(tt.input2, func(a int, b int) bool { return a == b })
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceForEach(t *testing.T) {
	tests := []struct {
		name  string
		input []int
	}{
		{
			name:  "case",
			input: newArray(100),
		},
		{
			name:  "case",
			input: newArray(123),
		},
		{
			name:  "case",
			input: newArray(1000),
		},
		{
			name:  "case",
			input: newArray(1234),
		},
		{
			name:  "case",
			input: newArray(10000),
		},
		{
			name:  "case",
			input: newArray(12345),
		},
		{
			name:  "nil",
			input: nil,
		},
		{
			name:  "empty",
			input: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).ForEach(func(i int, v int) { assert.Equal(t, tt.input[i], v) }).ToSlice()
			assert.Equal(t, tt.input, got)

			got = NewSliceByComparable(tt.input).ForEach(func(i int, v int) { assert.Equal(t, tt.input[i], v) }).ToSlice()
			assert.Equal(t, tt.input, got)

			got = NewSliceByOrdered(tt.input).ForEach(func(i int, v int) { assert.Equal(t, tt.input[i], v) }).ToSlice()
			assert.Equal(t, tt.input, got)

			got = NewSliceByMapping[int, int, int](tt.input).ForEach(func(i int, v int) { assert.Equal(t, tt.input[i], v) }).ToSlice()
			assert.Equal(t, tt.input, got)

			got = NewSlice(tt.input).Parallel(10).ForEach(func(i int, v int) { assert.Equal(t, tt.input[i], v) }).ToSlice()
			assert.Equal(t, tt.input, got)

			var count int64 = 0
			NewSlice(tt.input).Parallel(10).ForEach(func(i int, v int) { atomic.AddInt64(&count, 1) })
			assert.Equal(t, int64(len(tt.input)), count)
		})
	}
}

func TestSliceFindFunc(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(v int) bool
		want      int
	}{
		{
			name:      "case",
			input:     []int{1, 2, 1, 2, 1},
			predicate: func(v int) bool { return v == 1 },
			want:      0,
		},
		{
			name:      "case",
			input:     []int{1, 2, 1, 2, 1},
			predicate: func(v int) bool { return v == 2 },
			want:      1,
		},
		{
			name:      "case",
			input:     []int{1, 2},
			predicate: func(v int) bool { return v == 3 },
			want:      -1,
		},
		{
			name:      "empty",
			input:     []int{},
			predicate: func(v int) bool { return v == 1 },
			want:      -1,
		},
		{
			name:      "nil",
			input:     nil,
			predicate: func(v int) bool { return v == 1 },
			want:      -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).FindFunc(tt.predicate)
			assert.Equal(t, tt.want, got)

			got = NewSlice(tt.input).Parallel(4).FindFunc(tt.predicate)
			if got == -1 && tt.want != got {
				assert.Equal(t, tt.input[tt.want], tt.input[got])
			}

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
			input:     []string{"a", "b", "c"},
			predicate: func(v string) bool { return v != "b" },
			want:      []string{"a", "c"},
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
			got := NewSlice(tt.input).Filter(tt.predicate).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSlice(tt.input).Parallel(2).Filter(tt.predicate).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByComparable(tt.input).Filter(tt.predicate).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByComparable(tt.input).Parallel(2).Filter(tt.predicate).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).Filter(tt.predicate).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).Parallel(2).Filter(tt.predicate).ToSlice()
			assert.Equal(t, tt.want, got)
		})
	}

	tests1 := []struct {
		name      string
		input     []int
		predicate func(v int) bool
		want      int
	}{
		{
			name:      "match",
			input:     newArray(100),
			predicate: func(v int) bool { return v < 100 },
		},
		{
			name:      "match",
			input:     newArray(200),
			predicate: func(v int) bool { return v < 200 },
		},
		{
			name:      "match",
			input:     newArray(300),
			predicate: func(v int) bool { return v > 300 },
		},
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t,
				NewSliceByOrdered(tt.input).Parallel(10).Filter(tt.predicate).ToSlice(),
				NewSliceByOrdered(tt.input).Filter(tt.predicate).ToSlice())
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
			name:  "case",
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

func TestSliceInsert(t *testing.T) {
	tests := []struct {
		name   string
		input1 []int
		input2 []int
		input3 int
		want   []int
	}{
		{
			name:   "case",
			input1: []int{1, 2, 3},
			input2: []int{4, 5},
			input3: 1,
			want:   []int{1, 4, 5, 2, 3},
		},
		{
			name:   "case",
			input1: []int{1, 2, 3},
			input2: []int{4, 5},
			input3: 3,
			want:   []int{1, 2, 3, 4, 5},
		},
		{
			name:   "case",
			input1: []int{1, 2, 3},
			input2: []int{4, 5},
			input3: 5,
			want:   []int{1, 2, 3, 4, 5},
		},
		{
			name:   "empty",
			input1: []int{},
			input2: []int{4, 5},
			input3: 0,
			want:   []int{4, 5},
		},
		{
			name:   "nil",
			input1: []int{},
			input2: []int{4, 5},
			input3: 0,
			want:   []int{4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input1).Insert(tt.input3, tt.input2...).ToSlice()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceDelete(t *testing.T) {
	tests := []struct {
		name   string
		input1 []int
		input2 int
		input3 int
		want   []int
	}{
		{
			name:   "case",
			input1: []int{1, 2, 3},
			input2: 1,
			input3: 2,
			want:   []int{1, 3},
		},
		{
			name:   "case",
			input1: []int{1, 2, 3},
			input2: 0,
			input3: 1,
			want:   []int{2, 3},
		},
		{
			name:   "empty",
			input1: []int{},
			input2: 0,
			input3: 1,
			want:   []int{},
		},
		{
			name:   "nil",
			input1: []int{},
			input2: 0,
			input3: 1,
			want:   []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input1).Delete(tt.input2, tt.input3).ToSlice()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceIsSorted(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  bool
	}{
		{
			name:  "case",
			input: []int{1, 2, 1, 5},
			want:  false,
		},
		{
			name:  "case",
			input: []int{-1, -2, -1, -5},
			want:  false,
		},
		{
			name:  "case",
			input: []int{10, 11, 12, 13},
			want:  true,
		},
		{
			name:  "case",
			input: []int{-1, -2, -3, -4},
			want:  false,
		},
		{
			name:  "case",
			input: []int{-4, -3, -2, -1},
			want:  true,
		},
		{
			name:  "empty",
			input: []int{},
			want:  true,
		},
		{
			name:  "nil",
			input: nil,
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSlice(tt.input).IsSortedFunc(func(a, b int) bool { return a < b })
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
			name:  "case",
			input: []int{1, 2, 1},
			limit: 5,
			want:  []int{1, 2, 1},
		},
		{
			name:  "case",
			input: []int{1, 2, 1},
			limit: 2,
			want:  []int{1, 2},
		},
		{
			name:  "case",
			input: []int{1, 2, 1},
			limit: 3,
			want:  []int{1, 2, 1},
		},
		{
			name:  "case",
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
			got := NewSlice(tt.input).Limit(tt.limit).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByComparable(tt.input).Limit(tt.limit).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).Limit(tt.limit).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByMapping[int, int, int](tt.input).Limit(tt.limit).ToSlice()
			assert.Equal(t, tt.want, got)
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
			name:   "case",
			input:  []int{1, 2, 1},
			mapper: func(i int) int { return i * 2 },
			want:   []int{2, 4, 2},
		},
		{
			name:   "case",
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
			got := NewSlice(tt.input).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSlice(tt.input).Parallel(2).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByComparable(tt.input).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByComparable(tt.input).Parallel(2).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).Parallel(2).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)
		})
	}

	tests = []struct {
		name   string
		input  []int
		mapper func(int) int
		want   []int
	}{
		{
			name:   "case",
			input:  newArray(100),
			mapper: func(i int) int { return i * 2 },
		},
		{
			name:   "case",
			input:  newArray(200),
			mapper: func(i int) int { return i * 3 },
		},
		{
			name:   "case",
			input:  newArray(300),
			mapper: func(i int) int { return i * 4 },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t,
				NewSliceByOrdered(tt.input).Parallel(10).Map(tt.mapper).ToSlice(),
				NewSliceByOrdered(tt.input).Map(tt.mapper).ToSlice())
		})
	}
}

func TestSliceMaxFunc(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "case",
			input: []int{1, 2, 1, 5},
			want:  5,
		},
		{
			name:  "case",
			input: []int{-1, -2, -1, -5},
			want:  -1,
		},
		{
			name:  "case",
			input: []int{10, 2, 1, 5},
			want:  10,
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
			got := NewSlice(tt.input).MaxFunc(func(a, b int) bool { return a > b })
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceMinFunc(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "case",
			input: []int{1, 2, 1, 5},
			want:  1,
		},
		{
			name:  "case",
			input: []int{10, 2, 3, 1},
			want:  1,
		},
		{
			name:  "case",
			input: []int{-1, -2, -3, -1},
			want:  -3,
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
			got := NewSliceByOrdered(tt.input).MinFunc(func(a, b int) bool { return a < b })
			assert.Equal(t, tt.want, got)
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
			name:        "case",
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
			name:  "case",
			input: []int{1, 2, 1, 5},
			less:  func(a, b int) bool { return a > b },
			want:  []int{5, 2, 1, 1},
		},
		{
			name:  "case",
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
			got := NewSlice(tt.input).SortFunc(tt.less).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByComparable(tt.input).SortFunc(tt.less).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).SortFunc(tt.less).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByMapping[int, int, int](tt.input).SortFunc(tt.less).ToSlice()
			assert.Equal(t, tt.want, got)
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
			name:  "case",
			input: []int{1, 2, 1, 5},
			less:  func(a, b int) bool { return a > b },
			want:  []int{5, 2, 1, 1},
		},
		{
			name:  "case",
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
			got := NewSlice(tt.input).SortStableFunc(tt.less).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByComparable(tt.input).SortStableFunc(tt.less).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).SortStableFunc(tt.less).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByMapping[int, int, int](tt.input).SortStableFunc(tt.less).ToSlice()
			assert.Equal(t, tt.want, got)
		})
	}
}
