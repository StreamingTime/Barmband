package messaging

import (
	"errors"
	"fmt"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"strconv"
	"strings"
)

type Message interface {
}

type SetupMessage struct {
	BarmbandId barmband.BarmbandId
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

	return nil, UnknownMessageError
}

func parseSetupMessage(message string) (Message, error) {

	helloMessageFormat := fmt.Sprintf("%s %%s", SetupMessagePrefix)

	var bandId string
	_, err := fmt.Sscanf(message, helloMessageFormat, &bandId)

	if err != nil {
		return nil, err
	}
	if bandId == "" {
		return nil, EmptyBandId
	}
	return &SetupMessage{
		BarmbandId: barmband.BarmbandId(stringToBytes(bandId)),
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
