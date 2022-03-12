package stream

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPartition(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		goroutine int
		want1     []partition
		want2     int
	}{
		{
			name:      "case",
			input:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			goroutine: 2,
			want1: []partition{
				{
					low:  0,
					high: 6,
				},
				{
					low:  6,
					high: 13,
				},
			},
		},
		{
			name:      "case",
			input:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			goroutine: 3,
			want1: []partition{
				{
					low:  0,
					high: 4,
				},
				{
					low:  4,
					high: 8,
				},
				{
					low:  8,
					high: 13,
				},
			},
		},
		{
			name:      "case",
			input:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			goroutine: 5,
			want1: []partition{
				{
					low:  0,
					high: 2,
				},
				{
					low:  2,
					high: 4,
				},
				{
					low:  4,
					high: 7,
				},
				{
					low:  7,
					high: 10,
				},
				{
					low:  10,
					high: 13,
				},
			},
		},
		{
			name:      "case",
			input:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			goroutine: 10,
			want1: []partition{
				{
					low:  0,
					high: 1,
				},
				{
					low:  1,
					high: 2,
				},
				{
					low:  2,
					high: 3,
				},
				{
					low:  3,
					high: 4,
				},
				{
					low:  4,
					high: 5,
				},
				{
					low:  5,
					high: 6,
				},
				{
					low:  6,
					high: 7,
				},
				{
					low:  7,
					high: 9,
				},
				{
					low:  9,
					high: 11,
				},
				{
					low:  11,
					high: 13,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := partitionHandler(tt.input, tt.goroutine)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
