package cartridge

import (
	"github.com/kijimaD/goboy/pkg/types"
)

type MBC interface {
	Write(addr types.Word, value byte)
	Read(addr types.Word) byte
	switchROMBank(bank int)
	switchRAMBank(bank int)
}
