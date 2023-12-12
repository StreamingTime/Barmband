package bandcommand

import (
	"fmt"
	"log"
	"slices"
	"sync"

	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/messaging"
)

type BandCommand interface {
	HandleMessage(msg messaging.Message)
	HandleSetupMessage(msg *messaging.SetupMessage)
	HandlePairFoundMessage(message *messaging.PairFoundMessage)
	HandleAbortMessage(message *messaging.AbortMessage)
	HandleRequestPartnerMessage(message *messaging.RequestPartnerMessage)
	GetBand(id barmband.BarmbandId) *barmband.Barmband
}

type DefaultBandCommand struct {
	barmbandsMutex    sync.RWMutex
	barmbands         []barmband.Barmband
	messageHandler    func(bc BandCommand, msg messaging.Message)
	pairs             []barmband.Pair
	pairsMutex        sync.RWMutex
	pairFoundCallback func(pair barmband.Pair)
}

func New(pairFoundCallback func(pair barmband.Pair)) *DefaultBandCommand {
	return &DefaultBandCommand{
		barmbands:         make([]barmband.Barmband, 0),
		messageHandler:    defaultMessageHandler,
		pairFoundCallback: pairFoundCallback,
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
	case messaging.RequestPartnerMessage:
		fmt.Println("Got request partner message")
		bc.HandleRequestPartnerMessage(&msg)
	case *messaging.RequestPartnerMessage:
		fmt.Println("Got request partner message")
		bc.HandleRequestPartnerMessage(msg)

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

	bc.setWantsPair(message.BarmbandId, false)

	i := slices.IndexFunc(bc.pairs, func(pair barmband.Pair) bool {
		return pair.First == message.BarmbandId || pair.Second == message.BarmbandId
	})

	if i < 0 {
		return
	}

	bc.setWantsPair(bc.pairs[i].First, false)
	bc.setWantsPair(bc.pairs[i].Second, false)

	bc.pairs = append(bc.pairs[:i], bc.pairs[i+1:]...)
}

func (bc *DefaultBandCommand) HandleRequestPartnerMessage(message *messaging.RequestPartnerMessage) {
	if !bc.isRegistered(message.BarmbandId) {
		log.Printf("Band %s is not registered\n", message.BarmbandId)
		return
	}

	if bc.hasMatch(message.BarmbandId) {
		log.Printf("Band %s already has match\n", message.BarmbandId)
		return
	}

	partnerId := bc.findPartnerFor(message.BarmbandId)

	if partnerId == nil {
		bc.setWantsPair(message.BarmbandId, true)
		log.Printf("Band %s wants pair\n", message.BarmbandId)

		return
	}

	pair := barmband.NewPair(message.BarmbandId, *partnerId)

	log.Printf("Found pair %s %s\n", pair.First, pair.Second)

	bc.pairsMutex.Lock()
	bc.pairs = append(bc.pairs, pair)
	bc.pairsMutex.Unlock()

	bc.setWantsPair(pair.First, false)
	bc.setWantsPair(pair.Second, false)

	bc.pairFoundCallback(pair)
}

func (bc *DefaultBandCommand) isRegistered(barmbandId barmband.BarmbandId) bool {
	return bc.GetBand(barmbandId) != nil
}

func (bc *DefaultBandCommand) hasMatch(barmbandId barmband.BarmbandId) bool {
	return slices.ContainsFunc(bc.pairs, func(pair barmband.Pair) bool {
		return pair.First == barmbandId || pair.Second == barmbandId
	})
}

func (bc *DefaultBandCommand) setWantsPair(bandId barmband.BarmbandId, wantsPair bool) {
	band := bc.GetBand(bandId)
	if band == nil {
		return
	}
	band.WantsPair = wantsPair
}

func (bc *DefaultBandCommand) findPartnerFor(bandId barmband.BarmbandId) *barmband.BarmbandId {
	bc.barmbandsMutex.RLock()
	defer bc.barmbandsMutex.RUnlock()

	var wantsPairBands []barmband.Barmband

	for _, otherBand := range bc.barmbands {
		if otherBand.Id == bandId {
			continue
		}

		if otherBand.WantsPair {
			wantsPairBands = append(wantsPairBands, otherBand)
		}
	}

	if len(wantsPairBands) == 0 {
		return nil
	}

	// TODO: implement better matching algorithm

	return &wantsPairBands[0].Id
}
