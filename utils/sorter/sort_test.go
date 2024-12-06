package sorter

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6}
	Sort(items, func(i1, i2 int) bool {
		return i1 < i2 // ascending order
	})

	assert.True(t, slices.Equal(items, []int{1, 2, 3, 4, 5, 6}))

	Sort(items, func(i1, i2 int) bool {
		return i1 > i2 // descending order
	})

	assert.True(t, slices.Equal(items, []int{6, 5, 4, 3, 2, 1}))
}
