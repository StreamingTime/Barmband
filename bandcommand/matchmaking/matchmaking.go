package matchmaking

import (
	"log"
	"slices"
	"sync"

	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
)

type Matchmaker interface {
	StartMatchmaker(pairsChan chan barmband.Pair)
	AddBand(bandId barmband.BarmbandId)
	RemoveBand(bandId barmband.BarmbandId)
}

type DefaultMatchmaker struct {
	incomingBandsChan chan barmband.BarmbandId
	bandIdsLock       sync.RWMutex
	bandIds           []barmband.BarmbandId
}

func New() *DefaultMatchmaker {
	return &DefaultMatchmaker{
		incomingBandsChan: make(chan barmband.BarmbandId),
		bandIds:           make([]barmband.BarmbandId, 0),
	}
}

func (mm *DefaultMatchmaker) StartMatchmaker(pairsChan chan barmband.Pair) {

	go func() {
		for bandId := range mm.incomingBandsChan {
			mm.bandIdsLock.Lock()

			pairBandId := mm.FindPairFor(bandId)

			if pairBandId == nil {
				log.Printf("No pair found for %s\n", bandId)
				mm.bandIds = append(mm.bandIds, bandId)
			} else {
				log.Printf("Found pair for %s: %s\n", bandId, *pairBandId)
				pairsChan <- mm.MakePair(bandId, *pairBandId)
			}

			mm.bandIdsLock.Unlock()
		}
	}()
}

func (mm *DefaultMatchmaker) AddBand(bandId barmband.BarmbandId) {
	mm.incomingBandsChan <- bandId
}

func (mm *DefaultMatchmaker) RemoveBand(bandId barmband.BarmbandId) {
	mm.bandIdsLock.Lock()
	defer mm.bandIdsLock.Unlock()

	i := slices.IndexFunc(mm.bandIds, func(b barmband.BarmbandId) bool {
		return b == bandId
	})

	if i < 0 {
		return
	}

	mm.bandIds = append(mm.bandIds[:i], mm.bandIds[i+1:]...)
}

func (mm *DefaultMatchmaker) MakePair(a barmband.BarmbandId, b barmband.BarmbandId) barmband.Pair {
	mm.RemoveBand(a)
	mm.RemoveBand(b)

	return barmband.NewPair(a, b)
}

func (mm *DefaultMatchmaker) FindPairFor(bandId barmband.BarmbandId) *barmband.BarmbandId {
	if len(mm.bandIds) == 0 {
		return nil
	}

	return &mm.bandIds[0]
}
