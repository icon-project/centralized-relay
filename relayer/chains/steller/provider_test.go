package steller

import (
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	"github.com/stretchr/testify/assert"
)

func TestGetSeqBatches(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	type testCase struct {
		name                      string
		fromSeq, toSeq, batchSize uint64
		expected                  []types.LedgerSeqBatch
	}

	testCases := []testCase{
		{
			name:      "case-0",
			fromSeq:   1,
			toSeq:     1,
			batchSize: 1,
			expected: []types.LedgerSeqBatch{
				{FromSeq: 1, ToSeq: 1},
			},
		},
		{
			name:      "case-1",
			fromSeq:   1,
			toSeq:     2,
			batchSize: 1,
			expected: []types.LedgerSeqBatch{
				{FromSeq: 1, ToSeq: 1},
				{FromSeq: 2, ToSeq: 2},
			},
		},
		{
			name:      "case-2",
			fromSeq:   1,
			toSeq:     4,
			batchSize: 2,
			expected: []types.LedgerSeqBatch{
				{FromSeq: 1, ToSeq: 2},
				{FromSeq: 3, ToSeq: 4},
			},
		},
		{
			name:      "case-3",
			fromSeq:   1,
			toSeq:     4,
			batchSize: 3,
			expected: []types.LedgerSeqBatch{
				{FromSeq: 1, ToSeq: 3},
				{FromSeq: 4, ToSeq: 4},
			},
		},
		{
			name:      "case-4",
			fromSeq:   1,
			toSeq:     4,
			batchSize: 4,
			expected: []types.LedgerSeqBatch{
				{FromSeq: 1, ToSeq: 4},
			},
		},
		{
			name:      "case-5",
			fromSeq:   1,
			toSeq:     4,
			batchSize: 5,
			expected: []types.LedgerSeqBatch{
				{FromSeq: 1, ToSeq: 4},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(subTest *testing.T) {
			// batchStream := getLedgerSeqBatchStream(done, testCase.fromSeq, testCase.toSeq, testCase.batchSize)
			// batches := []types.LedgerSeqBatch{}
			// for batch := range batchStream {
			// 	batches = append(batches, batch)
			// }
			batches := getSeqBatches(testCase.fromSeq, testCase.toSeq, testCase.batchSize)
			assert.Equal(subTest, testCase.expected, batches)
		})
	}
}
