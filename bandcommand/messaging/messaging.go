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

type RequestPartnerMessage struct {
	BarmbandId barmband.BarmbandId
}

type messageParser func(string) (Message, error)

const SetupMessagePrefix = "Hello"
const PairFoundMessagePrefix = "Pair found"
const AbortMessagePrefix = "Abort"
const RequestPartnerPrefix = "Request partner"

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

	var bandIdS string
	_, err := fmt.Sscanf(message, helloMessageFormat, &bandIdS)

	if err != nil {
		return nil, err
	}
	if bandIdS == "" {
		return nil, EmptyBandId
	}

	bandId, err := barmband.IdFromString(bandIdS)
	if err != nil {
		return nil, err
	}
	return &SetupMessage{
		BarmbandId: bandId,
	}, nil
}

func parsePairFoundMessage(message string) (Message, error) {

	pairFoundMessageFormat := fmt.Sprintf("%s %%s %%s", PairFoundMessagePrefix)

	var firstBandIdS string
	var secondBandIdS string
	_, err := fmt.Sscanf(message, pairFoundMessageFormat, &firstBandIdS, &secondBandIdS)

	if err != nil {
		return nil, err
	}
	if firstBandIdS == "" || secondBandIdS == "" {
		return nil, EmptyBandId
	}

	firstBandId, err := barmband.IdFromString(firstBandIdS)
	if err != nil {
		return nil, err
	}

	secondBandId, err := barmband.IdFromString(secondBandIdS)
	if err != nil {
		return nil, err
	}

	return &PairFoundMessage{
		FirstBarmbandId:  firstBandId,
		SecondBarmbandId: secondBandId,
	}, nil
}

func parseAbortMessage(message string) (Message, error) {

	abortMessageFormat := fmt.Sprintf("%s %%s", AbortMessagePrefix)

	var bandIdS string
	_, err := fmt.Sscanf(message, abortMessageFormat, &bandIdS)

	if err != nil {
		return nil, err
	}
	if bandIdS == "" {
		return nil, EmptyBandId
	}

	bandId, err := barmband.IdFromString(bandIdS)
	if err != nil {
		return nil, err
	}
	return &AbortMessage{
		BarmbandId: bandId,
	}, nil
}

func parseRequestPartnerMessage(message string) (Message, error) {

	requestPartnerMessageFormat := fmt.Sprintf("%s %%s", RequestPartnerPrefix)

	var bandIdS string
	_, err := fmt.Sscanf(message, requestPartnerMessageFormat, &bandIdS)

	if err != nil {
		return nil, err
	}
	if bandIdS == "" {
		return nil, EmptyBandId
	}

	bandId, err := barmband.IdFromString(bandIdS)
	if err != nil {
		return nil, err
	}

	return &RequestPartnerMessage{
		BarmbandId: bandId,
	}, nil
}
