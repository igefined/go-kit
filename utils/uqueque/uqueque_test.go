package uqueque

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	t.Run("queue", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("new value 1")
		queue.Push("new value 2")
		queue.Push("new value 3")

		assert.NotNil(t, queue.tail)
		assert.Equal(t, queue.tail.value, "new value 3")

		pop := queue.Pop()
		assert.NotNil(t, pop)
		assert.Equal(t, pop, "new value 1")
	})
}
