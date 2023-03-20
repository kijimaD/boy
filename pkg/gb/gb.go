package gb

import (
	"time"

	"github.com/kijimaD/goboy/pkg/cpu"
	"github.com/kijimaD/goboy/pkg/gpu"
	"github.com/kijimaD/goboy/pkg/interfaces/window"
	"github.com/kijimaD/goboy/pkg/interrupt"
	"github.com/kijimaD/goboy/pkg/timer"
	"github.com/kijimaD/goboy/pkg/types"
)

// CyclesPerFrame is cpu clock num for 1frame.
const CyclesPerFrame = 70224

// GB is gameboy emulator struct
type GB struct {
	currentCycle uint
	cpu          *cpu.CPU
	gpu          *gpu.GPU
	timer        *timer.Timer
	irq          *interrupt.Interrupt
	win          window.Window
}

// NewGB is gb initializer
func NewGB(cpu *cpu.CPU, gpu *gpu.GPU, timer *timer.Timer, irq *interrupt.Interrupt, win window.Window) *GB {
	return &GB{
		currentCycle: 0,
		cpu:          cpu,
		gpu:          gpu,
		timer:        timer,
		irq:          irq,
		win:          win,
	}
}

// Start is
func (g *GB) Start() {
	t := time.NewTicker(16 * time.Millisecond)
	for {
		select {
		case <-t.C:
			g.win.Render(g.next())
		}
	}
	t.Stop()
}
func (g *GB) next() types.ImageData {
	for {
		var cycles uint
		if g.gpu.DMAStarted() {
			g.gpu.Transfer()
			// https://github.com/Gekkio/mooneye-gb/blob/master/docs/accuracy.markdown#how-many-cycles-does-oam-dma-take
			cycles = 162
		} else {
			cycles = g.cpu.Step()
		}
		g.gpu.Step(cycles * 4)
		if overflowed := g.timer.Update(cycles); overflowed {
			g.irq.SetIRQ(interrupt.TimerOverflowFlag)
		}
		g.currentCycle += cycles * 4
		if g.currentCycle >= CyclesPerFrame {
			g.win.PollKey()
			g.currentCycle -= CyclesPerFrame
			return g.gpu.GetImageData()
		}
	}
}
