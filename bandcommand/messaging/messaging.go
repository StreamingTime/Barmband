package messaging

import (
	"errors"
	"fmt"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"strings"
)

type Message interface {
}

type SetupMessage struct {
	BarmbandId barmband.BarmbandId
}

type PairFoundMessage struct {
	FirstBarmbandId  barmband.BarmbandId
	SecondBarmbandId barmband.BarmbandId
}

type AbortMessage struct {
	BarmbandId barmband.BarmbandId
}

type messageParser func(string) (Message, error)

const SetupMessagePrefix = "Hello"
const PairFoundMessagePrefix = "Pair found"
const AbortMessagePrefix = "Abort"

var UnknownMessageError = errors.New("unknown message")
var EmptyBandId = errors.New("EmptyBandId")

// messageParsers maps a message prefix to a messageParser
var messageParsers = map[string]messageParser{
	SetupMessagePrefix:     parseSetupMessage,
	PairFoundMessagePrefix: parsePairFoundMessage,
	AbortMessagePrefix:     parseAbortMessage,
}

// ParseMessage tries to convert a string into a Message using the parser configured in messageParsers
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
		BarmbandId: barmband.IdFromString(bandId),
	}, nil
}

func parsePairFoundMessage(message string) (Message, error) {

	pairFoundMessageFormat := fmt.Sprintf("%s %%s %%s", PairFoundMessagePrefix)

	var firstBandId string
	var secondBandId string
	_, err := fmt.Sscanf(message, pairFoundMessageFormat, &firstBandId, &secondBandId)

	if err != nil {
		return nil, err
	}
	if firstBandId == "" || secondBandId == "" {
		return nil, EmptyBandId
	}
	return &PairFoundMessage{
		FirstBarmbandId:  barmband.IdFromString(firstBandId),
		SecondBarmbandId: barmband.IdFromString(secondBandId),
	}, nil
}

func parseAbortMessage(message string) (Message, error) {

	abortMessageFormat := fmt.Sprintf("%s %%s", AbortMessagePrefix)

	var bandId string
	_, err := fmt.Sscanf(message, abortMessageFormat, &bandId)

	if err != nil {
		return nil, err
	}
	if bandId == "" {
		return nil, EmptyBandId
	}
	return &AbortMessage{
		BarmbandId: barmband.IdFromString(bandId),
	}, nil
}
