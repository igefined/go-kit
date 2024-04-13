package uslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsInt(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	assert.True(t, Contains(arr, 3))
	assert.False(t, Contains(arr, 6))
}

func TestContains_String(t *testing.T) {
	arr := []string{"apple", "banana", "orange", "grape"}
	assert.True(t, Contains(arr, "banana"))
	assert.False(t, Contains(arr, "kiwi"))
}

func TestContains_Float64(t *testing.T) {
	arr := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	assert.True(t, Contains(arr, 3.3))
	assert.False(t, Contains(arr, 6.6))
}
