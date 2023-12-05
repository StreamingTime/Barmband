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
		notSetupMessage := 1

		mockBc.EXPECT().HandleMessage(setupMsg).Times(1)

		defaultMessageHandler(mockBc, setupMsg)
		defaultMessageHandler(mockBc, notSetupMessage)
	})
}

func TestBandCommand_handleSetupMessage(t *testing.T) {

	bc := New()

	bandId := barmband.BarmbandId([]byte{1, 2, 3, 4})

	setupMsg := messaging.SetupMessage{BarmbandId: bandId}

	bc.handleSetupMessage(setupMsg)

	assert.Len(t, bc.barmbands, 1)
	assert.Equal(t, bandId, bc.barmbands[0].Id)

}
