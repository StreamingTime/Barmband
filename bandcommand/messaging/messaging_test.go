package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
)

var bandIdA = barmband.BarmbandId([]byte{0xAA, 0xAA, 0xAA, 0xAA})
var bandIdB = barmband.BarmbandId([]byte{0xBB, 0xBB, 0xBB, 0xBB})

func Test_ParseSetupMessage(t *testing.T) {

	t.Run("converts raw message to struct", func(t *testing.T) {

		rawMessage := "Hello AAAAAAAA"

		msg, err := parseSetupMessage(rawMessage)

		setupMsg := msg.(*SetupMessage)
		assert.Nil(t, err)
		assert.Equal(t, bandIdA, setupMsg.BarmbandId)
	})
}

func Test_ParsePairFoundMessage(t *testing.T) {

	t.Run("converts raw message to struct", func(t *testing.T) {

		rawMessage := "Pair found AAAAAAAA BBBBBBBB"

		msg, err := parsePairFoundMessage(rawMessage)

		setupMsg := msg.(*PairFoundMessage)
		assert.Nil(t, err)
		assert.Equal(t, bandIdA, setupMsg.FirstBarmbandId)
		assert.Equal(t, bandIdB, setupMsg.SecondBarmbandId)
	})
}

func Test_parseAbortMessage(t *testing.T) {
	t.Run("converts raw message to struct", func(t *testing.T) {

		rawMessage := "Abort AAAAAAAA"

		msg, err := parseAbortMessage(rawMessage)

		abortMessage := msg.(*AbortMessage)
		assert.Nil(t, err)
		assert.Equal(t, bandIdA, abortMessage.BarmbandId)
	})
}

func Test_parseRequestPartnerMessage(t *testing.T) {
	t.Run("converts raw message to struct", func(t *testing.T) {

		rawMessage := "Request partner AAAAAAAA"

		msg, err := parseRequestPartnerMessage(rawMessage)

		requestParnerMessage := msg.(*RequestPartnerMessage)
		assert.Nil(t, err)
		assert.Equal(t, bandIdA, requestParnerMessage.BarmbandId)
	})
}
