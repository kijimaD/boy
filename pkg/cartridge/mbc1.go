package cartridge

import (
	"github.com/kijimaD/goboy/pkg/ram"
	"github.com/kijimaD/goboy/pkg/rom"
	"github.com/kijimaD/goboy/pkg/types"
)

// MBC1 is Memory Bank Controller 1
// MBC1は2つの異なる最大メモリモードを持つ
// <16Mbit ROM/8KByte RAM> or <4Mbit ROM/32KByte RAM>
type MBC1 struct {
	rom             *rom.ROM
	ram             *ram.RAM
	selectedROMBank int
	selectedRAMBank int
	hasRAM          bool
	ramEnabled      bool
	hasBattery      bool
	memoryMode      MBC1MemoryMode
	RAMSize         int
}

// MBC1MemoryMode はMBC1のメモリ最大モード
// The MBC1 defaults to 16Mbit ROM/8KByte RAM mode
// MBC1のデフォルトは16Mbit ROM/8KByte RAM mode -- つまりROMを多く取る
// Writing a value (XXXXXXXS - X = Don't care, S = Memory model select) into 6000-7FFF area
// will select the memory model to use.
// S = 0 selects 16/8 mode. S = 1 selects 4/32 mode.

type MBC1MemoryMode = string

const (
	// ROM16mRAM8kMode is 4/32 memory mode
	// Writing a value (XXXXXXBB - X = Don't care, B = bank select bits) into 4000-5FFF area
	// will set the two most significant ROM address lines.
	// * NOTE: The Super Smart Card doesn't require this operation because it's RAM bank is ALWAYS enabled.
	// Include this operation anyway to allow your code to work with both
	ROM16mRAM8kMode MBC1MemoryMode = "ROM16M/RAM8K"
	// ROM4mRAM32kMode is 4/32 memory mode
	// Writing a value (XXXXXXBB - X = Don't care, B = bank select bits) into 4000-5FFF area
	// will select an appropriate RAM bank at A000-C000.
	// Before you can read or write to a RAM bank you have to enable it by writing a XXXX1010 into 0000-1FFF area*.
	// To disable RAM bank operations write any value but XXXX1010 into 0000-1FFF area.
	// Disabling a RAM bank probably protects that bank from false writes during power down of the GameBoy.
	// (NOTE: Nintendo suggests values $0A to enable and $00 to disable RAM bank!!)
	ROM4mRAM32kMode = "ROM4M/RAM32K"
)

func NewMBC1(buf []byte, ramSize int, hasBattery bool) *MBC1 {
	m := &MBC1{}
	m.memoryMode = ROM16mRAM8kMode
	m.hasBattery = hasBattery
	m.RAMSize = ramSize
	if ramSize > 0 {
		m.hasRAM = true
		m.ramEnabled = true
		m.selectedRAMBank = 0
		m.ram = ram.NewRAM(0x8000)
	}

	m.selectedROMBank = 0
	m.rom = rom.NewROM(buf)

	return m
}

func (m *MBC1) Write(addr types.Word, value byte) {
	switch {
	case addr >= 0x0000 && addr <= 0x1FFF:
		if m.memoryMode == ROM4mRAM32kMode && m.hasRAM {
			m.ramEnabled = value&0x0F == 0x0A
		}
		// Writing a value (XXXBBBBB - X = Don't cares, B = bank select bits) into 2000-3FFF area
		// will select an appropriate ROM bank at 4000-7FFF
		// Values of 0 and 1 do the same thing and point to ROM bank 1.
		// Rom bank 0 is not accessible from 4000-7FFF and can only be read from 0000-3FFF.
	case addr >= 0x2000 && addr <= 0x3FFF:
		m.switchROMBank(int(value & 0x1F))
	case addr >= 0x4000 && addr <= 0x5FFF:
		m.switchROMBank(int(value & 0x03))
	case addr >= 0x6000 && addr <= 0x7FFF:
		if value&0x01 == 0x00 {
			m.memoryMode = ROM16mRAM8kMode
		} else {
			m.memoryMode = ROM4mRAM32kMode
		}
	case addr >= 0xA000 && addr <= 0xBFFF:
		if m.hasRAM && m.ramEnabled {
			switch m.memoryMode {
			case ROM4mRAM32kMode:
				m.ram.Write(types.Word((int(addr)+m.selectedRAMBank*0x2000)-0xA000), value)
			case ROM16mRAM8kMode:
				m.ram.Write(types.Word((int(addr))-0xA000), value)
			}
		}
	}
}

func (m *MBC1) Read(addr types.Word) byte {
	if addr < 0x4000 {
		return m.rom.Read(addr)
	} else if addr >= 0x4000 && addr < 0x8000 {
		base := types.Word(0x0000)
		if m.selectedROMBank >= 2 {
			base = types.Word((m.selectedROMBank - 1) * 0x4000)
		}
		return m.rom.Read(base + addr)
	} else if addr >= 0xA000 && addr < 0xC000 {
		if m.hasRAM && m.ramEnabled {
			switch m.memoryMode {
			case ROM4mRAM32kMode:
				m.ram.Read(types.Word((int(addr) + m.selectedRAMBank*0x2000) - 0xA000))
			case ROM16mRAM8kMode:
				m.ram.Read(types.Word((int(addr)) - 0xA000))
			}
		}
	}
	return 0x00
}

func (m *MBC1) switchROMBank(bank int) {
	m.selectedROMBank = bank
}

func (m *MBC1) switchRAMBank(bank int) {
	m.selectedRAMBank = bank
}
