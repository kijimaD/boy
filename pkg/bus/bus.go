package bus

import (
	"github.com/kijimaD/goboy/pkg/interfaces/pad"
	"github.com/kijimaD/goboy/pkg/interrupt"

	"github.com/kijimaD/goboy/pkg/cartridge"
	"github.com/kijimaD/goboy/pkg/gpu"
	"github.com/kijimaD/goboy/pkg/interfaces/logger"
	"github.com/kijimaD/goboy/pkg/ram"
	"github.com/kijimaD/goboy/pkg/serial"
	"github.com/kijimaD/goboy/pkg/timer"
	"github.com/kijimaD/goboy/pkg/types"
	"github.com/kijimaD/goboy/pkg/utils"
)

const (
	// DMGStatusReg is DMA status register
	DMGStatusReg types.Word = 0xFF50
)

const (
	BANK_BEGIN             = 0x0000
	BANK_END               = 0x7FFF
	CARTRIDGE_HEADER_BEGIN = 0x0100
	VRAM_BEGIN             = 0x8000
	VRAM_END               = 0x9FFF
	EXT_RAM_BEGIN          = 0xA000
	EXT_RAM_END            = 0xBFFF
	WRAM_BEGIN             = 0xC000
	WRAM_END               = 0xDFFF
	ECHO_RAM_BEGIN         = 0xE000
	ECHO_RAM_END           = 0xFDFF
	OAM_BEGIN              = 0xFE00
	OAM_END                = 0xFE9F
	IO_REG_BEGIN           = 0xFF00
	IO_REG_END             = 0xFF7F
	JOYPAD                 = 0xFF00
	SERIAL_BEGIN           = 0xFF01
	TIMER_BEGIN            = 0xFF04
	TIMER_END              = 0xFF07
	GPU_BEGIN              = 0xFF40
	HRAM_BEGIN             = 0xFF80
	HRAM_END               = 0xFFFE
	IE_REG_ENABLE          = 0xFFFF
)

// Bus is gb bus
type Bus struct {
	logger    logger.Logger
	bootmode  bool
	cartridge *cartridge.Cartridge
	gpu       *gpu.GPU
	vRAM      *ram.RAM
	wRAM      *ram.RAM
	hRAM      *ram.RAM
	oamRAM    *ram.RAM
	timer     *timer.Timer
	irq       *interrupt.Interrupt
	pad       pad.Pad
}

/* --------------------------+
| Interrupt Enable Register  |
------------------------------  0xFFFF
| Internal RAM               |
------------------------------  0xFF80
| Empty but unusable for I/O |
------------------------------  0xFF4C
| I/O ports                  |
------------------------------  0xFF00
| Empty but unusable for I/O |
------------------------------  0xFEA0
| Sprite Attrib Memory (OAM) |
------------------------------  0xFE00
| Echo of 8kB Internal RAM   |
------------------------------  0xE000
| 8kB Internal RAM           |
------------------------------  0xC000
| 8kB switchable RAM bank    |
------------------------------  0xA000
| 8kB Video RAM              |
------------------------------  0x8000 --+
| 16kB switchable ROM bank   |           |
------------------------------  0x4000   | =  32kB Cartrigbe
| 16kB ROM bank #0           |           |
------------------------------  0x0000 --+   */

// NewBus is bus constructor
func NewBus(
	logger logger.Logger,
	cartridge *cartridge.Cartridge,
	gpu *gpu.GPU,
	vram *ram.RAM,
	wram *ram.RAM,
	hRAM *ram.RAM,
	oamRAM *ram.RAM,
	timer *timer.Timer,
	irq *interrupt.Interrupt,
	pad pad.Pad) *Bus {
	return &Bus{
		logger:    logger,
		bootmode:  true,
		cartridge: cartridge,
		gpu:       gpu,
		vRAM:      vram,
		wRAM:      wram,
		hRAM:      hRAM,
		oamRAM:    oamRAM,
		timer:     timer,
		irq:       irq,
		pad:       pad,
	}
}

// READBYTE is byte data reader from bus
// メモリマップ
func (b *Bus) ReadByte(addr types.Word) byte {
	switch {
	case addr >= BANK_BEGIN && addr <= BANK_END:
		if b.bootmode && addr < CARTRIDGE_HEADER_BEGIN {
			return BIOS[addr]
		}
		if addr == CARTRIDGE_HEADER_BEGIN {
			b.bootmode = false
		}
		return b.cartridge.ReadByte(addr)
	// Video RAM
	case addr >= VRAM_BEGIN && addr <= VRAM_END:
		return b.vRAM.Read(addr - VRAM_BEGIN)
	case addr >= EXT_RAM_BEGIN && addr <= EXT_RAM_END:
		return b.cartridge.ReadByte(addr)
	// Working RAM
	case addr >= WRAM_BEGIN && addr <= WRAM_END:
		return b.wRAM.Read(addr - WRAM_BEGIN)
	// Shadow
	case addr >= ECHO_RAM_BEGIN && addr <= ECHO_RAM_END:
		return b.wRAM.Read(addr - ECHO_RAM_BEGIN)
	// OAM
	case addr >= OAM_BEGIN && addr <= OAM_END:
		return b.oamRAM.Read(addr - OAM_BEGIN)
	// Pad
	case addr == JOYPAD:
		return b.pad.Read()
	// Timer
	case addr >= TIMER_BEGIN && addr <= TIMER_END:
		return b.timer.Read(addr - IO_REG_BEGIN)
	// IF
	case addr == 0xFF0F:
		return b.irq.Read(addr - IO_REG_BEGIN)
	// GPU
	case addr >= GPU_BEGIN && addr <= IO_REG_END:
		return b.gpu.Read(addr - GPU_BEGIN)
	// Zero page RAM
	case addr >= HRAM_BEGIN && addr <= HRAM_END:
		return b.hRAM.Read(addr - HRAM_BEGIN)
	// IE
	case addr == IE_REG_ENABLE:
		return b.irq.Read(addr - IO_REG_BEGIN)
	default:
		return 0
	}
	return 0
}

// ReadWord is word data reader from bus
func (b *Bus) ReadWord(addr types.Word) types.Word {
	l := b.ReadByte(addr)
	u := b.ReadByte(addr + 1)
	return utils.Bytes2Word(u, l)
}

// WriteByte is byte data writer to bus
func (b *Bus) WriteByte(addr types.Word, data byte) {
	switch {
	case addr >= BANK_BEGIN && addr <= BANK_END:
		b.cartridge.WriteByte(addr, data)
	// Video RAM
	case addr >= VRAM_BEGIN && addr <= VRAM_END:
		b.vRAM.Write(addr-VRAM_BEGIN, data)
	case addr >= EXT_RAM_BEGIN && addr <= EXT_RAM_END:
		b.cartridge.WriteByte(addr, data)
	// Working RAM
	case addr >= WRAM_BEGIN && addr <= WRAM_END:
		b.wRAM.Write(addr-WRAM_BEGIN, data)
	// Shadow
	case addr >= ECHO_RAM_BEGIN && addr <= ECHO_RAM_END:
		b.wRAM.Write(addr-ECHO_RAM_BEGIN, data)
	// OAM
	case addr >= OAM_BEGIN && addr <= OAM_END:
		b.oamRAM.Write(addr-OAM_BEGIN, data)
	// Pad
	case addr == IO_REG_BEGIN:
		b.pad.Write(data)
	// Serial
	case addr == SERIAL_BEGIN:
		serial.Send(data)
	// Timer
	case addr >= TIMER_BEGIN && addr <= TIMER_END:
		b.timer.Write(addr-IO_REG_BEGIN, data)
	// IF
	case addr == 0xFF0F:
		b.irq.Write(addr-IO_REG_BEGIN, data)
	// GPU
	case addr >= GPU_BEGIN && addr <= IO_REG_END:
		b.gpu.Write(addr-GPU_BEGIN, data)
	//Zero page RAM
	case addr >= HRAM_BEGIN && addr <= HRAM_END:
		b.hRAM.Write(addr-HRAM_BEGIN, data)
	// IE
	case addr == IE_REG_ENABLE:
		b.irq.Write(addr-IO_REG_BEGIN, data)
	default:
		// fmt.Printf("Error: You can not write 0x%X, this area is invalid or unimplemented area.\n", addr)
	}
}

// WriteWord is word data writer to bus
func (b *Bus) WriteWord(addr types.Word, data types.Word) {
	upper, lower := utils.Word2Bytes(data)
	b.WriteByte(addr, lower)
	b.WriteByte(addr+1, upper)
}

// BIOS is
var BIOS = []byte{ /* Removed*/ }
