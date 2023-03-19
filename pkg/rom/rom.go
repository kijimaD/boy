package rom

import "github.com/kijimaD/goboy/pkg/types"

// ROMから読み込んだ情報を保持
type ROM struct {
	data []byte
}

// construct
func NewROM(v []byte) *ROM {
	return &ROM{
		data: v,
	}
}

// 読み込み
func (r *ROM) Read(addr types.Word) byte {
	return r.data[addr]
}
