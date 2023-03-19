package cartridge

import (
	"github.com/kijimaD/goboy/pkg/rom"
	"github.com/kijimaD/goboy/pkg/types"
)

// MBC0 ROM ONLY
// This is a 32kb (256kb) ROM and occupies 0000-7FFF
type MBC0 struct {
	rom *rom.ROM
}

func NewMBC0(data []byte) *MBC0 {
	m := new(MBC0)
	m.rom = rom.NewROM(data)
	return m
}

func (m *MBC0) Write(addr types.Word, value byte) {
}

func (m *MBC0) Read(addr types.Word) byte {
	return m.rom.Read(addr)
}

func (m *MBC0) switchROMBank(bank int) {
	// ROM bankは1つなので
	// nop
}

func (m *MBC0) switchRAMBank(bank int) {
	// ROM bankは1つなので
	// nop
}
