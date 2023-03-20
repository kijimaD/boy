package main

import (
	"errors"
	"log"
	"os"

	"github.com/kijimaD/goboy/pkg/bus"
	"github.com/kijimaD/goboy/pkg/cartridge"
	"github.com/kijimaD/goboy/pkg/cpu"
	"github.com/kijimaD/goboy/pkg/gb"
	"github.com/kijimaD/goboy/pkg/gpu"
	"github.com/kijimaD/goboy/pkg/interrupt"
	"github.com/kijimaD/goboy/pkg/logger"
	"github.com/kijimaD/goboy/pkg/pad"
	"github.com/kijimaD/goboy/pkg/ram"
	"github.com/kijimaD/goboy/pkg/timer"
	"github.com/kijimaD/goboy/pkg/utils"
	"github.com/kijimaD/goboy/pkg/window"
)

// go run main.go roms/helloworld/hello.gb

func main() {
	level := "Debug"
	if os.Getenv("LEVEL") != "" {
		level = os.Getenv("LEVEL")
	}
	l := logger.NewLogger(logger.LogLevel(level))
	if len(os.Args) != 2 {
		log.Fatalf("ERROR: %v", errors.New("Please specify the ROM"))
	}
	file := os.Args[1]
	log.Println(file)
	buf, err := utils.LoadROM(file)
	if err != nil {
		log.Fatalf("ERROR: %v", errors.New("Failed to load ROM"))
	}
	cart, err := cartridge.NewCartridge(buf)
	if err != nil {
		log.Fatalf("ERROR: %v", errors.New("Failed to create cartridge"))
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
	win := window.NewWindow(pad)
	emu := gb.NewGB(cpu.NewCPU(l, b, irq), gpu, t, irq, win)
	win.Run(func() {
		win.Init()
		emu.Start()
	})
}
