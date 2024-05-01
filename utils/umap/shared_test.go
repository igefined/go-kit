package umap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharedMap(t *testing.T) {
	const n = 5

	contains := func(val string, values []string) bool {
		for _, v := range values {
			if val == v {
				return true
			}
		}
		return false
	}

	testKeys := []string{"alpha", "beta", "gamma"}

	sharedMap := NewSharedMap(n)
	sharedMap.Set(testKeys[0], 1)
	sharedMap.Set(testKeys[1], 2)
	sharedMap.Set(testKeys[2], 3)

	assert.Equal(t, 1, sharedMap.Get("alpha"))
	assert.Equal(t, 2, sharedMap.Get("beta"))
	assert.Equal(t, 3, sharedMap.Get("gamma"))

	keys := sharedMap.Keys()
	for _, k := range keys {
		assert.True(t, contains(k, testKeys))
	}
}
