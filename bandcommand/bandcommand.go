package bandcommand

import (
	"fmt"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/messaging"
	"slices"
)

type BandCommand interface {
	HandleMessage(msg messaging.Message)
}

type DefaultBandCommand struct {
	barmbands      []barmband.Barmband
	messageHandler func(bc BandCommand, msg messaging.Message)
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
		bc.HandleMessage(msg)
	}
}

func (bc *DefaultBandCommand) handleSetupMessage(setupMessage messaging.SetupMessage) {

	idAlreadyRegistered := slices.ContainsFunc(bc.barmbands, func(b barmband.Barmband) bool {
		return b.Id == setupMessage.BarmbandId
	})
	if idAlreadyRegistered {
		fmt.Printf("Band id %s is alredy registered\n", setupMessage.BarmbandId)
	} else {
		bc.barmbands = append(bc.barmbands, barmband.Barmband{
			Id: setupMessage.BarmbandId,
		})
	}
}
