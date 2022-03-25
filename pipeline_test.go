package stream

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipelines(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		predicate func(v string) bool
		mapper    func(v string) string
		want      []string
	}{
		{
			name:      "case",
			input:     []string{"a", "b", "c"},
			predicate: func(v string) bool { return v != "b" },
			mapper: func(v string) string {
				return v + "1"
			},
			want: []string{"a1", "c1"},
		},
		{
			name:      "case",
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
			got := NewSlice(tt.input).Filter(tt.predicate).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSlice(tt.input).Parallel(2).Filter(tt.predicate).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByComparable(tt.input).Filter(tt.predicate).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByComparable(tt.input).Parallel(2).Filter(tt.predicate).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).Filter(tt.predicate).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)

			got = NewSliceByOrdered(tt.input).Parallel(2).Filter(tt.predicate).Map(tt.mapper).ToSlice()
			assert.Equal(t, tt.want, got)
		})
	}
}
