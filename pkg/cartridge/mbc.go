package cartridge

import (
	"github.com/kijimaD/goboy/pkg/types"
)

// バンク切り替えによって利用可能なアドレス空間を拡張している
// カートリッジによってコントローラの違いがあるよう
type MBC interface {
	Write(addr types.Word, value byte)
	Read(addr types.Word) byte
	switchROMBank(bank int)
	switchRAMBank(bank int)
}
