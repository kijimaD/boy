package gpu

import (
	"testing"

	"github.com/kijimaD/goboy/pkg/constants"
	"github.com/kijimaD/goboy/pkg/interrupt"
	"github.com/kijimaD/goboy/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func setup() *GPU {
	g := NewGPU()
	irq := interrupt.NewInterrupt()
	g.Init(&mocks.MockBus{}, irq)
	return g
}

func TestLY(t *testing.T) {
	assert := assert.New(t)
	g := setup()
	for y := 0; y < int(constants.ScreenHeight+LCDVBlankHeight+10); y++ {
		if y == int(constants.ScreenHeight+LCDVBlankHeight) {
			assert.Equal(uint8(0x9a), g.Read(LY))
		} else {
			assert.Equal(byte(y%int(constants.ScreenHeight+LCDVBlankHeight)), g.Read(LY), y)
		}

		g.Step(CyclePerLine)
	}
}
