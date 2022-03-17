package stream

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestSliceMapping(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []string
		want2 string
	}{
		{
			name:  "case",
			input: []int{5, 1, 6, 2, 1},
			want:  []string{"mapping_1", "mapping_2", "mapping_1"},
			want2: "mapping_1/mapping_2/mapping_1/",
		},
		{
			name:  "empty",
			input: []int{},
			want:  []string{},
		},
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
	}
	f := func(v int) bool { return v < 5 }

	mapper := func(v int) string { return "mapping_" + strconv.Itoa(v) }
	reducer := func(r string, s string) string { return r + s + "/" }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := NewSliceByMapping[int, string, string](tt.input).Filter(f).Map(mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByMapping[int, string, string](tt.input).Filter(f).Parallel(2).Map(mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got1 := NewSliceByMapping[int, string, string](tt.input).Filter(f).Map(mapper).Reduce(reducer)
			assert.Equal(t, tt.want2, got1)

			got1 = NewSliceByMapping[int, string, string](tt.input).Filter(f).Parallel(2).Map(mapper).Reduce(reducer)
			assert.Equal(t, tt.want2, got1)

		})
	}

	tests = []struct {
		name  string
		input []int
		want  []string
		want2 string
	}{
		{
			name:  "case",
			input: newArray(100),
		},
		{
			name:  "case",
			input: newArray(200),
		},
		{
			name:  "case",
			input: newArray(300),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t,
				NewSliceByMapping[int, string, string](tt.input).Parallel(10).Map(mapper).ToSlice(),
				NewSliceByMapping[int, string, string](tt.input).Map(mapper).ToSlice())

			assert.Equal(t,
				NewSliceByMapping[int, string, string](tt.input).Parallel(10).Map(mapper).Reduce(reducer),
				NewSliceByMapping[int, string, string](tt.input).Map(mapper).Reduce(reducer))
		})
	}
}
