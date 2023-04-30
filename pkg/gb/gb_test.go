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
			"01-special",
			RomPathPrefix + "cpu_instrs/01-special.gb",
			1000,
		},
		{
			"02-interrupts",
			RomPathPrefix + "cpu_instrs/02-interrupts.gb",
			1000,
		},
		{
			"03-op sp,hl.gb",
			RomPathPrefix + "cpu_instrs/03-op sp,hl.gb",
			1000,
		},
		{
			"04-op r,imm.gb",
			RomPathPrefix + "cpu_instrs/04-op r,imm.gb",
			1000,
		},
		{
			"05-op rp.gb",
			RomPathPrefix + "cpu_instrs/05-op rp.gb",
			1000,
		},
		{
			"06-ld r,r.gb",
			RomPathPrefix + "cpu_instrs/06-ld r,r.gb",
			1000,
		},
		{
			"07-jr,jp,call,ret,rst.gb",
			RomPathPrefix + "cpu_instrs/07-jr,jp,call,ret,rst.gb",
			1000,
		},
		{
			"08-misc instrs.gb",
			RomPathPrefix + "cpu_instrs/08-misc instrs.gb",
			1000,
		},
		{
			"09-op r,r.gb",
			RomPathPrefix + "cpu_instrs/09-op r,r.gb",
			1000,
		},
		{
			"10-bit ops.gb",
			RomPathPrefix + "cpu_instrs/10-bit ops.gb",
			1000,
		},
		{
			"11-op a,(hl).gb",
			RomPathPrefix + "cpu_instrs/11-op a,(hl).gb",
			1000,
		},
		{
			"cpu_instr",
			RomPathPrefix + "cpu_instrs/cpu_instrs.gb",
			4000,
		},
		{
			"hello_world",
			RomPathPrefix + "helloworld/hello.gb",
			100,
		},
		{
			"genesis",
			RomPathPrefix + "genesis/Genesis1.gb",
			100,
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
