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
		want1     []partition[int]
		want2     int
	}{
		{
			name:      "case",
			input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			goroutine: 2,
			want1: []partition[int]{
				{
					slice:      []int{1, 2, 3, 4, 5, 6},
					startIndex: 0,
				},
				{
					slice:      []int{7, 8, 9, 10, 11, 12, 13},
					startIndex: 6,
				},
			},
		},
		{
			name:      "case",
			input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			goroutine: 3,
			want1: []partition[int]{
				{
					slice:      []int{1, 2, 3, 4},
					startIndex: 0,
				},
				{
					slice:      []int{5, 6, 7, 8},
					startIndex: 4,
				},
				{
					slice:      []int{9, 10, 11, 12, 13},
					startIndex: 8,
				},
			},
		},
		{
			name:      "case",
			input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			goroutine: 5,
			want1: []partition[int]{
				{
					slice:      []int{1, 2},
					startIndex: 0,
				},
				{
					slice:      []int{3, 4},
					startIndex: 2,
				},
				{
					slice:      []int{5, 6}, // expect 5,6,7
					startIndex: 4,           // expect 4
				},
				{
					slice:      []int{7, 8}, // expect 8,9,10
					startIndex: 6,           // expect 7
				},
				{
					slice:      []int{9, 10, 11, 12, 13}, // expect 11,12,13
					startIndex: 8,                        //expect 10
				},
			},
		},
		{
			name:      "case",
			input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			goroutine: 10,
			want1: []partition[int]{
				{
					slice:      []int{1},
					startIndex: 0,
				},
				{
					slice:      []int{2},
					startIndex: 1,
				},
				{
					slice:      []int{3},
					startIndex: 2,
				},
				{
					slice:      []int{4},
					startIndex: 3,
				},
				{
					slice:      []int{5},
					startIndex: 4,
				},
				{
					slice:      []int{6},
					startIndex: 5,
				},
				{
					slice:      []int{7},
					startIndex: 6,
				},
				{
					slice:      []int{8}, //expect 8,9
					startIndex: 7,        //expect 7
				},
				{
					slice:      []int{9}, //expect 10,11
					startIndex: 8,        //expect 9
				},
				{
					slice:      []int{10, 11, 12, 13}, //expect 12,13
					startIndex: 9,                     //expect 11
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
