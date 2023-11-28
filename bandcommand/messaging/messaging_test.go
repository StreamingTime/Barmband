package messaging

import (
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"testing"
)

func Test_ParseSetupMessage(t *testing.T) {

	t.Run("converts raw message to struct", func(t *testing.T) {

		rawMessage := "Hello 12345678"

		msg, err := parseSetupMessage(rawMessage)

		assert.Nil(t, err)
		assert.Len(t, msg.barmbandId, 4)
		assert.Equal(t, barmband.BarmbandId([]byte{0x12, 0x34, 0x56, 0x78}), msg.barmbandId)
	})
}

func Test_StringToBytes(t *testing.T) {

	input := "12345678"

	out := stringToBytes(input)

	assert.Equal(t, []byte{0x12, 0x34, 0x56, 0x78}, out)
}
