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

		setupMsg := msg.(*SetupMessage)
		assert.Nil(t, err)
		assert.Equal(t, barmband.BarmbandId([]byte{0x12, 0x34, 0x56, 0x78}), setupMsg.BarmbandId)
	})
}

func Test_ParsePairFoundMessage(t *testing.T) {

	t.Run("converts raw message to struct", func(t *testing.T) {

		rawMessage := "Pair found 12345678 11223344"

		msg, err := parsePairFoundMessage(rawMessage)

		setupMsg := msg.(*PairFoundMessage)
		assert.Nil(t, err)
		assert.Equal(t, barmband.BarmbandId([]byte{0x12, 0x34, 0x56, 0x78}), setupMsg.FirstBarmbandId)
		assert.Equal(t, barmband.BarmbandId([]byte{0x11, 0x22, 0x33, 0x44}), setupMsg.SecondBarmbandId)
	})
}

func Test_parseAbortMessage(t *testing.T) {
	t.Run("converts raw message to struct", func(t *testing.T) {

		rawMessage := "Abort AAAAAAAA"

		msg, err := parseAbortMessage(rawMessage)

		abortMessage := msg.(*AbortMessage)
		assert.Nil(t, err)
		assert.Equal(t, barmband.BarmbandId([]byte{0xAA, 0xAA, 0xAA, 0xAA}), abortMessage.BarmbandId)
	})
}
