// go:build units
package slice

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tLen = 1_000_000

func TestSearch(t *testing.T) {
	t.Run("success search test", func(t *testing.T) {
		nums := make([]int, tLen)
		for i := 0; i < tLen; i++ {
			nums[i] = i + 1
		}
		
		desired := rand.Intn(tLen)
		index := Search(nums, desired)

		assert.Equal(t, desired - 1, index)
	})

	t.Run("unsuccess search test", func(t *testing.T) {
		nums := make([]int, tLen)
		for i := 0; i < tLen; i++ {
			nums[i] = i + 1
		}
		
		desired := 1_234_5678
		index := Search(nums, desired)

		assert.Equal(t, -1, index)
	})
}