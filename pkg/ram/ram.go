package ram

import (
	"github.com/kijimaD/goboy/pkg/types"
)

// RAM is
type RAM struct {
	data []byte
}

// NewRAM is RAM constructor
func NewRAM(size int) *RAM {
	data := make([]byte, size)
	return &RAM{
		data: data,
	}
}

func (r *RAM) Read(addr types.Word) byte {
	return r.data[addr]
}

func (r *RAM) Write(addr types.Word, data byte) {
	r.data[addr] = data
}
