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

	data := []int{
		0b1111_1110, 0b1111_1100, // これで1行分。各ビットが1ピクセルに対応。2要素分をマスクして濃さを求める 例. 1x1 => 濃 1x0 => 薄
		0b1111_1110, 0b1111_1100,
		0b1111_1110, 0b1111_1100,
		0b1111_1110, 0b1111_1100,
		0b1111_1110, 0b1111_1100,
		0b1111_1110, 0b1111_1100,
		0b1111_1110, 0b0000_0000,
		0b0000_0000, 0b0000_0000,
	}

	// タイルは VRAM 0x8000 ~ 0x97FF
	addr := 0x8000
	for _, d := range data {
		g.bus.WriteByte(types.Word(addr), uint8(d))
		addr++
	}

	dataH := []int{
		0b0000_0000, 0b0000_0000,
		0b0100_0010, 0b0100_0010,
		0b0100_0010, 0b0100_0010,
		0b0111_1110, 0b0111_1110,
		0b0100_0010, 0b0100_0010,
		0b0100_0010, 0b0100_0010,
		0b0100_0010, 0b0100_0010,
		0b0000_0000, 0b0000_0000,
	}
	for _, d := range dataH {
		g.bus.WriteByte(types.Word(addr), uint8(d))
		addr++
	}

	dataE := []int{
		0b0000_0000, 0b0000_0000,
		0b0111_1110, 0b0111_1110,
		0b0100_0000, 0b0100_0000,
		0b0111_1110, 0b0111_1110,
		0b0100_0000, 0b0100_0000,
		0b0100_0000, 0b0100_0000,
		0b0111_1110, 0b0111_1110,
		0b0000_0000, 0b0000_0000,
	}
	for _, d := range dataE {
		g.bus.WriteByte(types.Word(addr), uint8(d))
		addr++
	}

	dataL := []int{
		0b0000_0000, 0b0000_0000,
		0b0100_0000, 0b0100_0000,
		0b0100_0000, 0b0100_0000,
		0b0100_0000, 0b0100_0000,
		0b0100_0000, 0b0100_0000,
		0b0100_0000, 0b0100_0000,
		0b0111_1110, 0b0111_1110,
		0b0000_0000, 0b0000_0000,
	}
	for _, d := range dataL {
		g.bus.WriteByte(types.Word(addr), uint8(d))
		addr++
	}

	dataO := []int{
		0b0000_0000, 0b0000_0000,
		0b0011_1100, 0b0011_1100,
		0b0100_0010, 0b0100_0010,
		0b0100_0010, 0b0100_0010,
		0b0100_0010, 0b0100_0010,
		0b0100_0010, 0b0100_0010,
		0b0011_1100, 0b0011_1100,
		0b0000_0000, 0b0000_0000,
	}
	for _, d := range dataO {
		g.bus.WriteByte(types.Word(addr), uint8(d))
		addr++
	}

	// タイルマップは0x9800 ~
	g.bus.WriteByte(types.Word(0x9800), uint8(0b0000_0001))
	g.bus.WriteByte(types.Word(0x9801), uint8(0b0000_0010))
	g.bus.WriteByte(types.Word(0x9802), uint8(0b0000_0011))
	g.bus.WriteByte(types.Word(0x9803), uint8(0b0000_0100))
	g.bus.WriteByte(types.Word(0x9804), uint8(0b0000_0100))

	g.bgPalette = 0b1110_0100

	for i := 0; i < 144; i++ {
		g.buildBGTile()
		g.ly = g.ly + 1
	}

	file, err := os.Create("../../test/unit/bgtile.png")
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

func TestBuildSprites(t *testing.T) {
	g := setup()

	tiledata := []int{
		0b000_0000, 0b0000_0000,
		0b0000_0000, 0b0000_0000,
		0b0000_0000, 0b0000_0000,
		0b0001_1000, 0b0000_0000,
		0b0001_1000, 0b0000_0000,
		0b0000_0000, 0b0000_0000,
		0b0000_0000, 0b0000_0000,
		0b0000_0000, 0b0000_0000,
	}

	// タイルは VRAM 0x8000 ~ 0x97FF
	addr := 0x8000
	for _, d := range tiledata {
		g.bus.WriteByte(types.Word(addr), uint8(d))
		addr++
	}

	spritedata := []int{
		0b0000_0000, 0b0000_0000,
		0b0000_0000, 0b0000_0000,
		0b0000_0000, 0b0110_0110,
		0b0000_0000, 0b0000_0000,
		0b0000_0000, 0b0000_0000,
		0b0000_0000, 0b0011_1100,
		0b0000_0000, 0b0000_0000,
		0b0000_0000, 0b0000_0000,
	}

	for _, d := range spritedata {
		g.bus.WriteByte(types.Word(addr), uint8(d))
		addr++
	}

	// タイルマップは0x9800 ~
	// 重なりを試すため
	g.bus.WriteByte(types.Word(0x9800), uint8(0b0000_0000)) // dot
	g.bus.WriteByte(types.Word(0x9801), uint8(0b0000_0001)) // face
	g.bus.WriteByte(types.Word(0x9802), uint8(0b0000_0010)) // empty

	// OAMは0xFE00 ~
	g.bus.WriteByte(types.Word(0xFE00), uint8(0b0010_1111)) // y
	g.bus.WriteByte(types.Word(0xFE01), uint8(0b0001_1111)) // x
	g.bus.WriteByte(types.Word(0xFE02), uint8(0b0000_0001)) // tilemap
	g.bus.WriteByte(types.Word(0xFE03), uint8(0b0000_0000)) // config

	g.bus.WriteByte(types.Word(0xFE04), uint8(0b0010_1111)) // y
	g.bus.WriteByte(types.Word(0xFE05), uint8(0b0010_1110)) // x
	g.bus.WriteByte(types.Word(0xFE06), uint8(0b0000_0001)) // tilemap
	g.bus.WriteByte(types.Word(0xFE07), uint8(0b0000_0000)) // config

	g.bgPalette = 0b1110_0100
	g.objPalette0 = 0b1110_0100
	g.objPalette1 = 0b1110_0100

	for i := 0; i < 144; i++ {
		g.buildBGTile()
		g.buildSprites()
		g.ly = g.ly + 1
	}

	file, err := os.Create("../../test/unit/sprite.png")
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

func TestGetTileID(t *testing.T) {
	assert := assert.New(t)
	g := setup()

	// タイルID書き込み
	g.bus.WriteByte(types.Word(0x9800), uint8(0b0000_1111))
	v := g.getTileID(0, uint(0), g.getWindowTilemapAddr())
	assert.Equal(int(0b0000_1111), v)

	// タイルID書き込み
	g.bus.WriteByte(types.Word(0x9801), uint8(0b0000_0111))
	v = g.getTileID(0, uint(1), g.getWindowTilemapAddr())
	assert.Equal(int(0b0000_0111), v)
}

func TestGetBGPalette(t *testing.T) {
	assert := assert.New(t)
	g := setup()

	// g.bgPaletteが指定されてないのですべて薄緑になる
	color := g.getBGPalette(3)
	assert.Equal(THIN_GREEN, color)
	color = g.getBGPalette(2)
	assert.Equal(THIN_GREEN, color)
	color = g.getBGPalette(1)
	assert.Equal(THIN_GREEN, color)
	color = g.getBGPalette(0)
	assert.Equal(THIN_GREEN, color)

	g.bgPalette = 0b1110_0100
	// 0b[11][10]_[01][00]
	color = g.getBGPalette(3)
	assert.Equal(BLACK_GREEN, color)
	color = g.getBGPalette(2)
	assert.Equal(DEEP_GREEN, color)
	color = g.getBGPalette(1)
	assert.Equal(MEDIUM_GREEN, color)
	color = g.getBGPalette(0)
	assert.Equal(THIN_GREEN, color)
}
