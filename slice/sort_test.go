// go:build units
package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	t.Run("success sort test", func(t *testing.T) {
		nums := []int{10, 87, 7, 90, 102, 65, 66, 32, 12, 71}

		Sort[int](nums)
		assert.Equal(t, []int{7, 10, 12, 32, 65, 66, 71, 87, 90, 102}, nums)
	})

	t.Run("success sort len = 2 test", func(t *testing.T) {
		nums := []int{87, 10}

		Sort[int](nums)
		assert.Equal(t, []int{10, 87}, nums)
	})

	t.Run("success sort len = 1 test", func(t *testing.T) {
		nums := []int{87}

		Sort[int](nums)
		assert.Equal(t, []int{87}, nums)
	})
}