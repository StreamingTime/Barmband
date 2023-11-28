package messaging

import (
	"errors"
	"fmt"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"strconv"
	"strings"
)

type Message = SetupMessage

type SetupMessage struct {
	barmbandId barmband.BarmbandId
}

type messageParser func(string) (Message, error)

const SetupMessagePrefix = "Hello"

var UnknownMessageError = errors.New("unknown message")
var EmptyBandId = errors.New("EmptyBandId")

var messageParsers = map[string]messageParser{
	SetupMessagePrefix: parseSetupMessage,
}

func ParseMessage(message string) (Message, error) {

	for prefix, parser := range messageParsers {
		if strings.HasPrefix(message, prefix) {
			return parser(message)
		}
	}

	return Message{}, UnknownMessageError
}

func parseSetupMessage(message string) (SetupMessage, error) {

	helloMessageFormat := fmt.Sprintf("%s %%s", SetupMessagePrefix)

	var bandId string
	_, err := fmt.Sscanf(message, helloMessageFormat, &bandId)

	if err != nil {
		return SetupMessage{}, err
	}
	if bandId == "" {
		return SetupMessage{}, EmptyBandId
	}
	return SetupMessage{
		barmbandId: barmband.BarmbandId(stringToBytes(bandId)),
	}, nil
}

// stringToBytes converts  the string "1234" to the byte slice []byte{0x12, 0x34}
func stringToBytes(s string) []byte {
	bytes := make([]byte, 0, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		num, _ := strconv.ParseUint(s[i:i+2], 16, 8)
		bytes = append(bytes, byte(num))
	}
	fmt.Printf("%X\n", bytes)
	return bytes
}
