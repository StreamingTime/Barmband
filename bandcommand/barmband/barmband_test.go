package barmband

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdFromString(t *testing.T) {

	input := "12345678"

	out := IdFromString(input)

	assert.Equal(t, BarmbandId([]byte{0x12, 0x34, 0x56, 0x78}), out)

}
