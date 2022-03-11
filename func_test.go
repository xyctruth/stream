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
		want1     [][]int
		want2     int
	}{
		{
			name:      "case",
			input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			goroutine: 2,
			want1:     [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12, 13}},
			want2:     6,
		},
		{
			name:      "case",
			input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			goroutine: 3,
			want1:     [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12, 13}},
			want2:     4,
		},
		{
			name:      "case",
			input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			goroutine: 5,
			want1:     [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10, 11, 12, 13}},
			//want1:     [][]int{{1, 2}, {3, 4}, {5, 6, 7}, {8, 9, 10}, {11, 12, 13}},
			want2: 2,
		},
		{
			name:      "case",
			input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			goroutine: 10,
			want1:     [][]int{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10, 11, 12, 13}},
			//want1:     [][]int{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8, 9}, {10, 11}, {12, 13}},
			want2: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := partition(tt.input, tt.goroutine)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)

		})
	}
}
