package rom

import (
	"github.com/kijimaD/goboy/pkg/types"
)

// ROM
type ROM struct {
	data []byte
}

// NewROM is ROM constructor
func NewROM(v []byte) *ROM {
	return &ROM{
		data: v,
	}
}

func (r *ROM) Read(addr types.Word) byte {
	return r.data[addr]
}
