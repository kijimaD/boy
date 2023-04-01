package gpu

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/kijimaD/goboy/pkg/constants"
	"github.com/kijimaD/goboy/pkg/interrupt"
	"github.com/kijimaD/goboy/pkg/mocks"
	"github.com/kijimaD/goboy/pkg/types"

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

func TestBuildBGTile(t *testing.T) {
	g := setup()
	// タイルは VRAM 8000 ~ 97FF
	data := []int{
		0b1111_1110, 0b1111_1100, // 各ビットで1列に対応。2要素分をマスクして濃さを求める 例. 1x1 => 濃 1x0 => 薄
		0b1111_1110, 0b1111_1100,
		0b1111_1110, 0b1111_1100,
		0b1111_1110, 0b1111_1100,
		0b1111_1110, 0b1111_1100,
		0b1111_1110, 0b1111_1100,
		0b1111_1110, 0b0000_0000,
		0b0000_0000, 0b0000_0000,
	}

	addr := 0x8000
	for _, d := range data {
		g.bus.WriteByte(types.Word(addr), uint8(d))
		addr++
	}

	g.bgPalette = 0b1110_0100

	for i := 0; i < 8; i++ {
		g.buildBGTile()
		g.ly = g.ly + 1
	}

	file, err := os.Create("../../test/unit/block.png")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, constants.ScreenWidth, constants.ScreenHeight))
	set(img, g.imageData)
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func set(img *image.RGBA, imageData types.ImageData) {
	rect := img.Rect
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			img.Set(x, rect.Max.Y-y, imageData[y*rect.Max.X+x])
		}
	}
}
