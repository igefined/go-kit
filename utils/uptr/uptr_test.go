package uptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	intValue := 10
	intPtr := Ptr(intValue)
	assert.NotNil(t, intPtr)
	assert.Equal(t, intValue, *intPtr)

	strValue := "test"
	strPtr := Ptr(strValue)
	assert.NotNil(t, strPtr)
	assert.Equal(t, strValue, *strPtr)

	type TestStruct struct {
		Field string
	}
	structValue := TestStruct{Field: "example"}
	structPtr := Ptr(structValue)
	assert.NotNil(t, structPtr)
	assert.Equal(t, structValue, *structPtr)
}
