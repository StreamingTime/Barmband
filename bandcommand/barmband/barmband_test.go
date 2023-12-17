package barmband

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdFromString(t *testing.T) {

	t.Run("converts string to BarmbandId", func(t *testing.T) {

		input := "12345678"

		out, err := IdFromString(input)

		assert.Nil(t, err)
		assert.Equal(t, BarmbandId([]byte{0x12, 0x34, 0x56, 0x78}), out)
	})

	t.Run("fails when len(s) is not 8", func(t *testing.T) {

		inputs := []string{
			"aa", "aaaaaaaaaaaaaaaa",
		}

		for _, s := range inputs {

			t.Run(fmt.Sprintf("s is %s", s), func(t *testing.T) {
				_, err := IdFromString(s)

				assert.NotNil(t, err)
			})

		}
	})
}

func TestIdToString(t *testing.T) {

	id := BarmbandId([]byte{0x12, 0x34, 0x56, 0x78})

	assert.Equal(t, "12345678", IdToString(id))
}
