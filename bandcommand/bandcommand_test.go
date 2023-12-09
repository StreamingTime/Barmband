package bandcommand

import (
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/_mocks/mock_bandcommand"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/messaging"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_defaultMessageHandler(t *testing.T) {

	t.Run("routes messages correctly", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		mockBc := mock_bandcommand.NewMockBandCommand(mockCtrl)

		bandId := barmband.BarmbandId([]byte{1, 2, 3, 4})

		setupMsg := messaging.SetupMessage{BarmbandId: bandId}
		setupMsg2 := messaging.SetupMessage{BarmbandId: barmband.BarmbandId([]byte{5, 5, 5, 5})}

		abortMsg := messaging.AbortMessage{BarmbandId: bandId}

		notSetupMessage := 1

		pairFoundMessage := messaging.PairFoundMessage{
			FirstBarmbandId:  barmband.BarmbandId([]byte{0xa, 0xa, 0xa, 0xa}),
			SecondBarmbandId: barmband.BarmbandId([]byte{0xb, 0xb, 0xb, 0xb}),
		}
		mockBc.EXPECT().HandleSetupMessage(gomock.Any()).Times(2)
		mockBc.EXPECT().HandlePairFoundMessage(&pairFoundMessage).Times(2)
		mockBc.EXPECT().HandleAbortMessage(&abortMsg).Times(2)

		defaultMessageHandler(mockBc, setupMsg)
		defaultMessageHandler(mockBc, &setupMsg2)
		defaultMessageHandler(mockBc, notSetupMessage)
		defaultMessageHandler(mockBc, pairFoundMessage)
		defaultMessageHandler(mockBc, &pairFoundMessage)

		defaultMessageHandler(mockBc, abortMsg)
		defaultMessageHandler(mockBc, &abortMsg)
	})
}

func TestBandCommand_handleSetupMessage(t *testing.T) {

	bc := New()

	bandId := barmband.BarmbandId([]byte{1, 2, 3, 4})

	setupMsg := messaging.SetupMessage{BarmbandId: bandId}

	bc.HandleSetupMessage(&setupMsg)

	assert.Len(t, bc.barmbands, 1)
	assert.Equal(t, bandId, bc.barmbands[0].Id)

}

func TestBandCommand_HandlePairFoundMessage(t *testing.T) {

	t.Run("handles existing match with different band order", func(t *testing.T) {

		bandA := barmband.BarmbandId([]byte{1, 2, 3, 4})
		bandB := barmband.BarmbandId([]byte{5, 5, 5, 5})
		bandC := barmband.BarmbandId([]byte{12, 12, 12, 12})

		type testCase struct {
			name   string
			first  barmband.BarmbandId
			second barmband.BarmbandId
		}
		testCases := []testCase{
			{
				name:   "A,B",
				first:  bandA,
				second: bandB,
			},
			{
				name:   "B,A",
				first:  bandB,
				second: bandA,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				bc := New()

				bc.barmbands = []barmband.Barmband{
					barmband.NewBarmband(bandA),
					barmband.NewBarmband(bandB),
					barmband.NewBarmband(bandC),
				}

				bc.pairs = append(bc.pairs, barmband.NewPair(tc.first, tc.second))

				pairFoundMsg := messaging.PairFoundMessage{FirstBarmbandId: tc.first, SecondBarmbandId: tc.second}

				bc.HandlePairFoundMessage(&pairFoundMsg)

				assert.Len(t, bc.pairs, 0)

				assert.Equal(t, 1, bc.GetBand(tc.first).FoundPairs)
				assert.Equal(t, 1, bc.GetBand(tc.second).FoundPairs)
				assert.Equal(t, 0, bc.GetBand(bandC).FoundPairs)
			})

		}

	})

	t.Run("ignores non existent matches", func(t *testing.T) {

		bandA := barmband.BarmbandId([]byte{1, 2, 3, 4})
		bandB := barmband.BarmbandId([]byte{5, 5, 5, 5})
		bandC := barmband.BarmbandId([]byte{12, 12, 12, 12})

		bc := New()

		bc.barmbands = []barmband.Barmband{
			barmband.NewBarmband(bandA),
			barmband.NewBarmband(bandB),
			barmband.NewBarmband(bandC),
		}

		// Match: A and B
		bc.pairs = append(bc.pairs, barmband.NewPair(bandA, bandB))

		// Message reports match A and C
		pairFoundMsg := messaging.PairFoundMessage{FirstBarmbandId: bandA, SecondBarmbandId: bandC}

		bc.HandlePairFoundMessage(&pairFoundMsg)

		assert.Len(t, bc.pairs, 1)

		for _, b := range bc.barmbands {
			assert.Equal(t, 0, b.FoundPairs)
		}

	})

}

func TestBandCommand_HandleAbortMessage(t *testing.T) {

	t.Run("handles existing match with different band order", func(t *testing.T) {

		bandA := barmband.BarmbandId([]byte{1, 2, 3, 4})
		bandB := barmband.BarmbandId([]byte{5, 5, 5, 5})
		bandC := barmband.BarmbandId([]byte{12, 12, 12, 12})

		type testCase struct {
			name   string
			bandId barmband.BarmbandId
		}
		testCases := []testCase{
			{
				name:   "A",
				bandId: bandA,
			},
			{
				name:   "B",
				bandId: bandB,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				bc := New()

				bc.barmbands = []barmband.Barmband{
					barmband.NewBarmband(bandA),
					barmband.NewBarmband(bandB),
					barmband.NewBarmband(bandC),
				}

				bc.pairs = []barmband.Pair{
					{
						First:  bandA,
						Second: bandB,
					},
				}

				abortMsg := messaging.AbortMessage{BarmbandId: tc.bandId}

				bc.HandleAbortMessage(&abortMsg)

				assert.Len(t, bc.pairs, 0)
			})

		}

	})

	t.Run("ignores abort when band has no match", func(t *testing.T) {

		bandA := barmband.BarmbandId([]byte{1, 2, 3, 4})
		bandB := barmband.BarmbandId([]byte{5, 5, 5, 5})
		bandC := barmband.BarmbandId([]byte{12, 12, 12, 12})

		bc := New()

		bc.barmbands = []barmband.Barmband{
			barmband.NewBarmband(bandA),
			barmband.NewBarmband(bandB),
			barmband.NewBarmband(bandC),
		}

		// Match: A and B
		bc.pairs = append(bc.pairs, barmband.NewPair(bandA, bandB))

		// Message reports match A and C
		abortMsg := messaging.AbortMessage{BarmbandId: bandC}

		bc.HandleAbortMessage(&abortMsg)

		assert.Len(t, bc.pairs, 1)

	})

}

func TestDefaultBandCommand_GetBand(t *testing.T) {
	bandA := barmband.BarmbandId([]byte{1, 2, 3, 4})
	bandB := barmband.BarmbandId([]byte{5, 5, 5, 5})
	bandC := barmband.BarmbandId([]byte{12, 12, 12, 12})

	bc := New()

	bc.barmbands = []barmband.Barmband{
		barmband.NewBarmband(bandA),
		barmband.NewBarmband(bandB),
	}

	assert.Equal(t, bandB, bc.GetBand(bandB).Id)
	assert.Nil(t, bc.GetBand(bandC))
}
