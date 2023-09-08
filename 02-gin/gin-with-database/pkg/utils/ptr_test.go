package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPtr(t *testing.T) {
	i := 12
	res := ToPtr[int](i)
	assert.Equal(t, &i, res)
}
