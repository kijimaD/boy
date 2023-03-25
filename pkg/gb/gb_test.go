package gb

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/kijimaD/goboy/pkg/bus"
	"github.com/kijimaD/goboy/pkg/cartridge"
	"github.com/kijimaD/goboy/pkg/constants"
	"github.com/kijimaD/goboy/pkg/cpu"
	"github.com/kijimaD/goboy/pkg/gpu"
	"github.com/kijimaD/goboy/pkg/interfaces/window"
	"github.com/kijimaD/goboy/pkg/interrupt"
	"github.com/kijimaD/goboy/pkg/logger"
	"github.com/kijimaD/goboy/pkg/pad"
	"github.com/kijimaD/goboy/pkg/ram"
	"github.com/kijimaD/goboy/pkg/timer"
	"github.com/kijimaD/goboy/pkg/types"
	"github.com/kijimaD/goboy/pkg/utils"
)

const (
	RomPathPrefix   = "../../roms/"
	ImagePathPrefix = "../../test/actual/"
)

// MockWindow is
type mockWindow struct {
	window.Window
}

func (m mockWindow) PollKey() {
}

func setup(file string) *GB {
	l := logger.NewLogger(logger.LogLevel("DEBUG"))
	buf, err := utils.LoadROM(file)
	if err != nil {
		panic(err)
	}
	cart, err := cartridge.NewCartridge(buf)
	if err != nil {
		panic(err)
	}

	vRAM := ram.NewRAM(0x2000)
	wRAM := ram.NewRAM(0x2000)
	hRAM := ram.NewRAM(0x80)
	oamRAM := ram.NewRAM(0xA0)
	gpu := gpu.NewGPU()
	t := timer.NewTimer()
	pad := pad.NewPad()
	irq := interrupt.NewInterrupt()
	b := bus.NewBus(l, cart, gpu, vRAM, wRAM, hRAM, oamRAM, t, irq, pad)
	gpu.Init(b, irq)
	win := mockWindow{}
	emu := NewGB(cpu.NewCPU(l, b, irq), gpu, t, irq, win)
	return emu
}

func set(img *image.RGBA, imageData types.ImageData) {
	rect := img.Rect
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			img.Set(x, rect.Max.Y-y, imageData[y*rect.Max.X+x])
		}
	}
}

func skipFrame(emu *GB, n int) types.ImageData {
	var image types.ImageData
	for i := 0; i < n; i++ {
		image = emu.next()
	}
	return image
}

func TestROMs(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		frame int
	}{
		{
			"hello_world",
			RomPathPrefix + "helloworld/hello.gb",
			100,
		},
		{
			"cpu_instr",
			RomPathPrefix + "cpu_instrs/cpu_instrs.gb",
			4000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emu := setup(tt.path)
			imageData := skipFrame(emu, tt.frame)
			file, err := os.Create(ImagePathPrefix + tt.name + ".png")
			defer file.Close()
			if err != nil {
				panic(err)
			}
			img := image.NewRGBA(image.Rect(0, 0, constants.ScreenWidth, constants.ScreenHeight))
			set(img, imageData)
			if err := png.Encode(file, img); err != nil {
				panic(err)
			}
		})
	}
}
