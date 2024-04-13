package validator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/igefined/go-kit/utils/uptr"
)

func TestNewValidator(t *testing.T) {
	validator, err := NewValidator()
	assert.NoError(t, err)
	assert.NotNil(t, validator)
}

func TestDateTimeOrNil(t *testing.T) {
	type TestData struct {
		Time *time.Time `validate:"omitempty,dateTimeOrNil"`
	}

	validate, err := NewValidator()
	assert.NoError(t, err)

	nilTimeData := TestData{Time: nil}

	err = validate.Struct(nilTimeData)
	assert.NoError(t, err)

	futureTimeData := TestData{Time: uptr.Ptr(time.Now().Add(time.Hour))}
	err = validate.Struct(futureTimeData)
	assert.NoError(t, err)
}

func TestXor(t *testing.T) {
	tCases := []struct {
		x      bool
		y      bool
		result bool
	}{
		{
			x:      true,
			y:      true,
			result: false,
		},
		{
			x:      true,
			y:      false,
			result: true,
		},
		{
			x:      false,
			y:      true,
			result: true,
		},
		{
			x:      false,
			y:      false,
			result: false,
		},
	}

	for _, c := range tCases {
		t.Run("test xor", func(t *testing.T) {
			assert.Equal(t, c.result, xor(c.x, c.y))
		})
	}
}
