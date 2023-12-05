package bandcommand

import (
	"fmt"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/messaging"
	"slices"
)

type BandCommand interface {
	HandleMessage(msg messaging.Message)
	HandleSetupMessage(msg *messaging.SetupMessage)
	HandlePairFoundMessage(message *messaging.PairFoundMessage)
}

type DefaultBandCommand struct {
	barmbands      []barmband.Barmband
	messageHandler func(bc BandCommand, msg messaging.Message)
	pairs          []barmband.Pair
}

func New() *DefaultBandCommand {
	return &DefaultBandCommand{
		barmbands:      make([]barmband.Barmband, 0),
		messageHandler: defaultMessageHandler,
	}
}

func (bc *DefaultBandCommand) HandleMessage(msg messaging.Message) {
	bc.messageHandler(bc, msg)
}

func defaultMessageHandler(bc BandCommand, msg messaging.Message) {

	switch msg.(type) {
	case messaging.SetupMessage:
		fmt.Println("got setup message")
		var setupMessage messaging.SetupMessage = msg.(messaging.SetupMessage)
		bc.HandleSetupMessage(&setupMessage)
	case *messaging.SetupMessage:
		fmt.Println("got setup message")
		var setupMessage *messaging.SetupMessage = msg.(*messaging.SetupMessage)
		bc.HandleSetupMessage(setupMessage)
	case messaging.PairFoundMessage:
		fmt.Println("Got pair found message")
		var pairFoundMessage messaging.PairFoundMessage = msg.(messaging.PairFoundMessage)
		bc.HandlePairFoundMessage(&pairFoundMessage)
	case *messaging.PairFoundMessage:
		fmt.Println("Got pair found message")
		var pairFoundMessage *messaging.PairFoundMessage = msg.(*messaging.PairFoundMessage)
		bc.HandlePairFoundMessage(pairFoundMessage)

	default:
		fmt.Printf("Unknown message: %T\n", msg)
	}
}

func (bc *DefaultBandCommand) HandleSetupMessage(setupMessage *messaging.SetupMessage) {

	idAlreadyRegistered := slices.ContainsFunc(bc.barmbands, func(b barmband.Barmband) bool {
		return b.Id == setupMessage.BarmbandId
	})
	if idAlreadyRegistered {
		fmt.Printf("Band id %s is alredy registered\n", setupMessage.BarmbandId)
	} else {
		bc.barmbands = append(bc.barmbands, barmband.Barmband{
			Id: setupMessage.BarmbandId,
		})
		fmt.Printf("registered band %v\n", setupMessage.BarmbandId)
	}
}

func (bc *DefaultBandCommand) HandlePairFoundMessage(message *messaging.PairFoundMessage) {
	i := slices.IndexFunc(bc.pairs, func(p barmband.Pair) bool {
		return (p.First == message.FirstBarmbandId && p.Second == message.SecondBarmbandId) || (p.First == message.SecondBarmbandId && p.Second == message.FirstBarmbandId)
	})

	if i < 0 {
		fmt.Println("No pair found")
		return
	}

	bc.pairs = append(bc.pairs[:i], bc.pairs[i+1:]...)
}
