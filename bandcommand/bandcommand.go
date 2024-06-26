package bandcommand

import (
	"log"
	"slices"
	"sync"

	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/color"
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
		bc.HandleSetupMessage(&msg)
	case *messaging.SetupMessage:
		bc.HandleSetupMessage(msg)
	case messaging.PairFoundMessage:
		bc.HandlePairFoundMessage(&msg)
	case *messaging.PairFoundMessage:
		bc.HandlePairFoundMessage(msg)
	case messaging.AbortMessage:
		bc.HandleAbortMessage(&msg)
	case *messaging.AbortMessage:
		bc.HandleAbortMessage(msg)
	case messaging.RequestPartnerMessage:
		bc.HandleRequestPartnerMessage(&msg)
	case *messaging.RequestPartnerMessage:
		bc.HandleRequestPartnerMessage(msg)

	default:
		log.Printf("Unknown message: %T\n", msg)
	}
}

func (bc *DefaultBandCommand) HandleSetupMessage(setupMessage *messaging.SetupMessage) {

	bc.barmbandsMutex.Lock()
	defer bc.barmbandsMutex.Unlock()

	idAlreadyRegistered := slices.ContainsFunc(bc.barmbands, func(b barmband.Barmband) bool {
		return b.Id == setupMessage.BarmbandId
	})
	if idAlreadyRegistered {
		log.Printf("Band id %s is already registered\n", barmband.IdToString(setupMessage.BarmbandId))
	} else {
		bc.barmbands = append(bc.barmbands, barmband.Barmband{
			Id: setupMessage.BarmbandId,
		})
		log.Printf("Registered band %s\n", barmband.IdToString(setupMessage.BarmbandId))
	}
}

func (bc *DefaultBandCommand) HandlePairFoundMessage(message *messaging.PairFoundMessage) {
	bc.pairsMutex.Lock()
	defer bc.pairsMutex.Unlock()

	i := slices.IndexFunc(bc.pairs, func(p barmband.Pair) bool {
		return (p.First == message.FirstBarmbandId && p.Second == message.SecondBarmbandId) || (p.First == message.SecondBarmbandId && p.Second == message.FirstBarmbandId)
	})

	if i < 0 {
		log.Println("Pair could not be found")
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
		log.Printf("Band %s is not registered\n", barmband.IdToString(message.BarmbandId))
		return
	}

	if bc.hasMatch(message.BarmbandId) {
		log.Printf("Band %s already has match\n", barmband.IdToString(message.BarmbandId))
		return
	}

	partnerId := bc.findPartnerFor(message.BarmbandId)

	if partnerId == nil {
		bc.setWantsPair(message.BarmbandId, true)
		log.Printf("Band %s wants pair\n", barmband.IdToString(message.BarmbandId))

		return
	}

	c, err := bc.GetRandomColor()

	if err != nil {
		log.Printf("Error getting random color: %s\n", err)
		return
	}

	pair := barmband.NewPair(message.BarmbandId, *partnerId, c)

	log.Printf("Found pair %s %s\n", barmband.IdToString(message.BarmbandId), barmband.IdToString(message.BarmbandId))

	bc.pairsMutex.Lock()
	bc.pairs = append(bc.pairs, pair)
	bc.pairsMutex.Unlock()

	bc.setWantsPair(pair.First, false)
	bc.setWantsPair(pair.Second, false)

	bc.pairFoundCallback(pair)
}

func (bc *DefaultBandCommand) GetRandomColor() (string, error) {
	return color.GetRandomColor(bc.GetUsedColors())
}

func (bc *DefaultBandCommand) GetUsedColors() []string {
	bc.pairsMutex.RLock()

	var colors []string

	for _, pair := range bc.pairs {
		colors = append(colors, pair.Color)
	}

	bc.pairsMutex.RUnlock()

	return colors
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
