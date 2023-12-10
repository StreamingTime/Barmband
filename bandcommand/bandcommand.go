package bandcommand

import (
	"fmt"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/messaging"
	"slices"
	"sync"
)

type BandCommand interface {
	HandleMessage(msg messaging.Message)
	HandleSetupMessage(msg *messaging.SetupMessage)
	HandlePairFoundMessage(message *messaging.PairFoundMessage)
	HandleAbortMessage(message *messaging.AbortMessage)
	GetBand(id barmband.BarmbandId) *barmband.Barmband
}

type DefaultBandCommand struct {
	barmbandsMutex sync.RWMutex
	barmbands      []barmband.Barmband
	messageHandler func(bc BandCommand, msg messaging.Message)
	pairs          []barmband.Pair
	pairsMutex     sync.RWMutex
}

func New() *DefaultBandCommand {
	return &DefaultBandCommand{
		barmbands:      make([]barmband.Barmband, 0),
		messageHandler: defaultMessageHandler,
	}
}

// GetBand returns a pointer to the band with the given id, or nil if there is no band with this id
func (bc *DefaultBandCommand) GetBand(id barmband.BarmbandId) *barmband.Barmband {
	bc.barmbandsMutex.RLock()
	defer bc.barmbandsMutex.RUnlock()

	i := slices.IndexFunc(bc.barmbands, func(b barmband.Barmband) bool {
		return b.Id == id
	})

	if i < 0 {
		return nil
	}

	return &bc.barmbands[i]
}

func (bc *DefaultBandCommand) HandleMessage(msg messaging.Message) {
	bc.messageHandler(bc, msg)
}

func defaultMessageHandler(bc BandCommand, msg messaging.Message) {

	switch msg := msg.(type) {
	case messaging.SetupMessage:
		fmt.Println("got setup message")
		bc.HandleSetupMessage(&msg)
	case *messaging.SetupMessage:
		fmt.Println("got setup message")
		bc.HandleSetupMessage(msg)
	case messaging.PairFoundMessage:
		fmt.Println("Got pair found message")
		bc.HandlePairFoundMessage(&msg)
	case *messaging.PairFoundMessage:
		fmt.Println("Got pair found message")
		bc.HandlePairFoundMessage(msg)
	case messaging.AbortMessage:
		fmt.Println("Got abort message")
		bc.HandleAbortMessage(&msg)
	case *messaging.AbortMessage:
		fmt.Println("Got abort message")
		bc.HandleAbortMessage(msg)

	default:
		fmt.Printf("Unknown message: %T\n", msg)
	}
}

func (bc *DefaultBandCommand) HandleSetupMessage(setupMessage *messaging.SetupMessage) {

	bc.barmbandsMutex.Lock()
	defer bc.barmbandsMutex.Unlock()

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
	bc.pairsMutex.Lock()
	defer bc.pairsMutex.Unlock()

	i := slices.IndexFunc(bc.pairs, func(p barmband.Pair) bool {
		return (p.First == message.FirstBarmbandId && p.Second == message.SecondBarmbandId) || (p.First == message.SecondBarmbandId && p.Second == message.FirstBarmbandId)
	})

	if i < 0 {
		fmt.Println("No pair found")
		return
	}

	bandA := bc.GetBand(message.FirstBarmbandId)
	bandB := bc.GetBand(message.SecondBarmbandId)

	bandA.FoundPairs++
	bandB.FoundPairs++

	bc.pairs = append(bc.pairs[:i], bc.pairs[i+1:]...)
}

func (bc *DefaultBandCommand) HandleAbortMessage(message *messaging.AbortMessage) {

	bc.pairsMutex.Lock()
	defer bc.pairsMutex.Unlock()

	i := slices.IndexFunc(bc.pairs, func(pair barmband.Pair) bool {
		return pair.First == message.BarmbandId || pair.Second == message.BarmbandId
	})

	if i < 0 {
		return
	}

	bc.pairs = append(bc.pairs[:i], bc.pairs[i+1:]...)
}
