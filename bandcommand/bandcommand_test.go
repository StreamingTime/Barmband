package bandcommand

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/_mocks/mock_bandcommand"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/messaging"
	"go.uber.org/mock/gomock"
)

var bandIdA = barmband.BarmbandId([]byte{0xAA, 0xAA, 0xAA, 0xAA})
var bandIdB = barmband.BarmbandId([]byte{0xBB, 0xBB, 0xBB, 0xBB})
var bandIdC = barmband.BarmbandId([]byte{0xCC, 0xCC, 0xCC, 0xCC})
var defaultColor = "red"

func Test_defaultMessageHandler(t *testing.T) {

	t.Run("routes messages correctly", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		mockBc := mock_bandcommand.NewMockBandCommand(mockCtrl)

		setupMsg := messaging.SetupMessage{BarmbandId: bandIdA}
		setupMsg2 := messaging.SetupMessage{BarmbandId: bandIdB}

		abortMsg := messaging.AbortMessage{BarmbandId: bandIdA}

		pairFoundMessage := messaging.PairFoundMessage{
			FirstBarmbandId:  bandIdA,
			SecondBarmbandId: bandIdB,
		}

		requestPartnerMessage := messaging.RequestPartnerMessage{BarmbandId: bandIdA}

		mockBc.EXPECT().HandleSetupMessage(gomock.Any()).Times(2)
		mockBc.EXPECT().HandlePairFoundMessage(&pairFoundMessage).Times(2)
		mockBc.EXPECT().HandleAbortMessage(&abortMsg).Times(2)
		mockBc.EXPECT().HandleRequestPartnerMessage(&requestPartnerMessage).Times(2)

		defaultMessageHandler(mockBc, setupMsg)
		defaultMessageHandler(mockBc, &setupMsg2)
		defaultMessageHandler(mockBc, pairFoundMessage)
		defaultMessageHandler(mockBc, &pairFoundMessage)

		defaultMessageHandler(mockBc, abortMsg)
		defaultMessageHandler(mockBc, &abortMsg)

		defaultMessageHandler(mockBc, requestPartnerMessage)
		defaultMessageHandler(mockBc, &requestPartnerMessage)
	})
}

func TestBandCommand_handleSetupMessage(t *testing.T) {

	bc := New(func(pair barmband.Pair) {
		t.Fatalf("pairFoundCallback should not be called")
	})

	setupMsg := messaging.SetupMessage{BarmbandId: bandIdA}

	bc.HandleSetupMessage(&setupMsg)

	assert.Len(t, bc.barmbands, 1)
	assert.Equal(t, bandIdA, bc.barmbands[0].Id)

}

func TestBandCommand_HandlePairFoundMessage(t *testing.T) {

	t.Run("handles existing match with different band order", func(t *testing.T) {

		type testCase struct {
			name   string
			first  barmband.BarmbandId
			second barmband.BarmbandId
		}
		testCases := []testCase{
			{
				name:   "A,B",
				first:  bandIdA,
				second: bandIdB,
			},
			{
				name:   "B,A",
				first:  bandIdB,
				second: bandIdA,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				bc := New(func(pair barmband.Pair) {
					t.Fatalf("pairFoundCallback should not be called")
				})

				bc.barmbands = []barmband.Barmband{
					barmband.NewBarmband(bandIdA),
					barmband.NewBarmband(bandIdB),
					barmband.NewBarmband(bandIdC),
				}

				bc.pairs = append(bc.pairs, barmband.NewPair(tc.first, tc.second, defaultColor))

				pairFoundMsg := messaging.PairFoundMessage{FirstBarmbandId: tc.first, SecondBarmbandId: tc.second}

				bc.HandlePairFoundMessage(&pairFoundMsg)

				assert.Len(t, bc.pairs, 0)

				assert.Equal(t, 1, bc.GetBand(tc.first).FoundPairs)
				assert.Equal(t, 1, bc.GetBand(tc.second).FoundPairs)
				assert.Equal(t, 0, bc.GetBand(bandIdC).FoundPairs)
			})

		}

	})

	t.Run("ignores non existent matches", func(t *testing.T) {

		bc := New(func(pair barmband.Pair) {
			t.Fatalf("pairFoundCallback should not be called")
		})

		bc.barmbands = []barmband.Barmband{
			barmband.NewBarmband(bandIdA),
			barmband.NewBarmband(bandIdB),
			barmband.NewBarmband(bandIdC),
		}

		// Match: A and B
		bc.pairs = append(bc.pairs, barmband.NewPair(bandIdA, bandIdB, defaultColor))

		// Message reports match A and C
		pairFoundMsg := messaging.PairFoundMessage{FirstBarmbandId: bandIdA, SecondBarmbandId: bandIdC}

		bc.HandlePairFoundMessage(&pairFoundMsg)

		assert.Len(t, bc.pairs, 1)

		for _, b := range bc.barmbands {
			assert.Equal(t, 0, b.FoundPairs)
		}

	})

}

func TestBandCommand_HandleAbortMessage(t *testing.T) {

	t.Run("handles existing match with different band order", func(t *testing.T) {

		type testCase struct {
			name   string
			bandId barmband.BarmbandId
		}
		testCases := []testCase{
			{
				name:   "A",
				bandId: bandIdA,
			},
			{
				name:   "B",
				bandId: bandIdB,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				bc := New(func(pair barmband.Pair) {
					t.Fatalf("pairFoundCallback should not be called")
				})

				bc.barmbands = []barmband.Barmband{
					barmband.NewBarmband(bandIdA),
					barmband.NewBarmband(bandIdB),
					barmband.NewBarmband(bandIdC),
				}

				bc.pairs = []barmband.Pair{
					{
						First:  bandIdA,
						Second: bandIdB,
					},
				}

				abortMsg := messaging.AbortMessage{BarmbandId: tc.bandId}

				bc.HandleAbortMessage(&abortMsg)

				assert.Len(t, bc.pairs, 0)
				assert.False(t, bc.GetBand(bandIdA).WantsPair)
				assert.False(t, bc.GetBand(bandIdB).WantsPair)
			})

		}

	})

	t.Run("ignores abort when band has no match", func(t *testing.T) {

		bc := New(func(pair barmband.Pair) {
			t.Fatalf("pairFoundCallback should not be called")
		})

		bc.barmbands = []barmband.Barmband{
			barmband.NewBarmband(bandIdA),
			barmband.NewBarmband(bandIdB),
			barmband.NewBarmband(bandIdC),
		}

		// Match: A and B
		bc.pairs = append(bc.pairs, barmband.NewPair(bandIdA, bandIdB, defaultColor))

		// Message reports match A and C
		abortMsg := messaging.AbortMessage{BarmbandId: bandIdC}

		bc.HandleAbortMessage(&abortMsg)

		assert.Len(t, bc.pairs, 1)

	})

}

func TestBandCommand_HandleRequestPartnerMessage(t *testing.T) {

	t.Run("is ignored when there is a pair with this band", func(t *testing.T) {

		type testCase struct {
			name   string
			bandId barmband.BarmbandId
		}
		testCases := []testCase{
			{
				name:   "A",
				bandId: bandIdA,
			},
			{
				name:   "B",
				bandId: bandIdB,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				bc := New(func(pair barmband.Pair) {
					t.Fatalf("pairFoundCallback should not be called")
				})

				bc.barmbands = []barmband.Barmband{
					barmband.NewBarmband(bandIdA),
					barmband.NewBarmband(bandIdB),
					barmband.NewBarmband(bandIdC),
				}

				bc.pairs = []barmband.Pair{
					{
						First:  bandIdA,
						Second: bandIdB,
						Color:  defaultColor,
					},
				}

				requestPair := messaging.RequestPartnerMessage{BarmbandId: tc.bandId}

				bc.HandleRequestPartnerMessage(&requestPair)

				assert.Len(t, bc.pairs, 1)
				assert.False(t, bc.GetBand(tc.bandId).WantsPair)
			})

		}

	})

	t.Run("ignores request when band is not registered", func(t *testing.T) {

		bc := New(func(pair barmband.Pair) {
			t.Fatalf("pairFoundCallback should not be called")
		})

		msg := messaging.RequestPartnerMessage{BarmbandId: bandIdA}

		bc.HandleRequestPartnerMessage(&msg)

		assert.Len(t, bc.pairs, 0)
	})

	t.Run("ignores request when band is already waiting for a partner", func(t *testing.T) {

		bc := New(func(pair barmband.Pair) {
			t.Fatalf("pairFoundCallback should not be called")
		})

		bandA := barmband.NewBarmband(bandIdA)
		bandA.WantsPair = true

		bandB := barmband.NewBarmband(bandIdA)
		bandC := barmband.NewBarmband(bandIdA)

		bc.barmbands = []barmband.Barmband{
			bandA,
			bandB,
			bandC,
		}

		msg := messaging.RequestPartnerMessage{BarmbandId: bandIdA}

		bc.HandleRequestPartnerMessage(&msg)

		assert.Len(t, bc.pairs, 0)
		assert.True(t, bandA.WantsPair)
		assert.False(t, bandB.WantsPair)
		assert.False(t, bandC.WantsPair)
	})

	t.Run("requests partner when no match can be created", func(t *testing.T) {

		bc := New(func(pair barmband.Pair) {
			t.Fatalf("pairFoundCallback should not be called")
		})

		bc.barmbands = []barmband.Barmband{
			barmband.NewBarmband(bandIdA),
			barmband.NewBarmband(bandIdB),
		}

		msg := messaging.RequestPartnerMessage{BarmbandId: bandIdA}

		bc.HandleRequestPartnerMessage(&msg)

		assert.Len(t, bc.pairs, 0)
		assert.True(t, bc.GetBand(bandIdA).WantsPair)
		assert.False(t, bc.GetBand(bandIdB).WantsPair)
	})

	t.Run("creates match when possible", func(t *testing.T) {

		bc := New(func(pair barmband.Pair) {
			assert.Equal(t, bandIdA, pair.Second)
			assert.Equal(t, bandIdB, pair.First)
		})

		bc.barmbands = []barmband.Barmband{
			barmband.NewBarmband(bandIdA),
			barmband.NewBarmband(bandIdB),
		}

		msg1 := messaging.RequestPartnerMessage{BarmbandId: bandIdA}
		bc.HandleRequestPartnerMessage(&msg1)

		msg2 := messaging.RequestPartnerMessage{BarmbandId: bandIdB}
		bc.HandleRequestPartnerMessage(&msg2)

		assert.Len(t, bc.pairs, 1)
		assert.False(t, bc.GetBand(bandIdA).WantsPair)
		assert.False(t, bc.GetBand(bandIdB).WantsPair)
	})

}

func TestDefaultBandCommand_GetBand(t *testing.T) {

	bc := New(func(pair barmband.Pair) {
		t.Fatalf("pairFoundCallback should not be called")
	})

	bc.barmbands = []barmband.Barmband{
		barmband.NewBarmband(bandIdA),
		barmband.NewBarmband(bandIdB),
	}

	assert.Equal(t, bandIdB, bc.GetBand(bandIdB).Id)
	assert.Nil(t, bc.GetBand(bandIdC))
}

func TestDefaultBandCommand_isRegistered(t *testing.T) {

	bc := New(func(pair barmband.Pair) {
		t.Fatalf("pairFoundCallback should not be called")
	})

	bc.barmbands = append(bc.barmbands, barmband.NewBarmband(bandIdA), barmband.NewBarmband(bandIdB))

	assert.True(t, bc.isRegistered(bandIdA))
	assert.True(t, bc.isRegistered(bandIdB))
	assert.False(t, bc.isRegistered(bandIdC))

}

func TestDefaultBandCommand_hasMatch(t *testing.T) {

	bc := New(func(pair barmband.Pair) {
		t.Fatalf("pairFoundCallback should not be called")
	})

	bc.pairs = append(bc.pairs, barmband.NewPair(bandIdA, bandIdB, defaultColor))

	assert.True(t, bc.hasMatch(bandIdA))
	assert.True(t, bc.hasMatch(bandIdB))
	assert.False(t, bc.hasMatch(bandIdC))

}

func TestDefaultBandCommand_setWantsPair(t *testing.T) {

	bc := New(func(pair barmband.Pair) {
		t.Fatalf("pairFoundCallback should not be called")
	})

	bc.barmbands = append(bc.barmbands, barmband.NewBarmband(bandIdA), barmband.NewBarmband(bandIdB))

	bc.setWantsPair(bandIdA, true)
	bc.setWantsPair(bandIdB, false)

	assert.True(t, bc.GetBand(bandIdA).WantsPair)
	assert.False(t, bc.GetBand(bandIdB).WantsPair)
}

func TestDefaultBandCommand_findPartnerFor(t *testing.T) {

	bc := New(func(pair barmband.Pair) {
		t.Fatalf("pairFoundCallback should not be called")
	})

	bandA := barmband.NewBarmband(bandIdA)
	bandA.WantsPair = true

	bandB := barmband.NewBarmband(bandIdB)
	bandB.WantsPair = false

	bc.barmbands = append(bc.barmbands, bandA, bandB)

	assert.Equal(t, bandIdA, *bc.findPartnerFor(bandIdB))
	assert.Nil(t, bc.findPartnerFor(bandIdA))
}

func TestGetUsedColors(t *testing.T) {
	bc := &DefaultBandCommand{
		pairs: []barmband.Pair{
			{Color: "red"},
			{Color: "blue"},
			{Color: "green"},
		},
	}

	expected := []string{"red", "blue", "green"}
	actual := bc.GetUsedColors()

	assert.Equal(t, expected, actual)
}

func TestGetRandomColor(t *testing.T) {
	bc := &DefaultBandCommand{
		pairs: []barmband.Pair{
			{Color: "red"},
			{Color: "blue"},
			{Color: "green"},
		},
	}

	color, err := bc.GetRandomColor()
	assert.Nil(t, err)
	assert.NotContains(t, bc.GetUsedColors(), color)
}
