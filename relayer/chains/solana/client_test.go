package solana

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPercentileItem(t *testing.T) {
	tests := []struct {
		name       string
		fees       []uint64
		percentile uint64
		want       uint64
	}{
		{
			name:       "basic sorted array with 50th percentile",
			fees:       []uint64{1, 2, 3, 4, 5},
			percentile: 50,
			want:       3,
		},
		{
			name:       "single element array",
			fees:       []uint64{42},
			percentile: 50,
			want:       42,
		},
		{
			name:       "unsorted array",
			fees:       []uint64{5, 2, 8, 1, 9},
			percentile: 50,
			want:       5,
		},
		{
			name:       "array with duplicates",
			fees:       []uint64{1, 2, 2, 2, 3},
			percentile: 50,
			want:       2,
		},
		{
			name:       "0th percentile returns max",
			fees:       []uint64{1, 2, 3, 4, 5},
			percentile: 0,
			want:       5,
		},
		{
			name:       "100th percentile returns max",
			fees:       []uint64{1, 2, 3, 4, 5},
			percentile: 100,
			want:       5,
		},
		{
			name:       "25th percentile",
			fees:       []uint64{1, 2, 3, 4, 5, 6, 7, 8},
			percentile: 25,
			want:       2,
		},
		{
			name:       "75th percentile",
			fees:       []uint64{1, 2, 3, 4, 5, 6, 7, 8},
			percentile: 75,
			want:       6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feesCopy := make([]uint64, len(tt.fees))
			copy(feesCopy, tt.fees)

			got, err := getPercentileItem(feesCopy, tt.percentile)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
