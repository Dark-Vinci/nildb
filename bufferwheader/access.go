package bufferwheader

import (
	"unsafe"

	"github.com/dark-vinci/nildb/utils"
)

func (bwh *BufferWithHeader[H]) Header() *H {
	return bwh.header
}

func (bwh *BufferWithHeader[H]) Content() []byte {
	return bwh.content
}

func (bwh *BufferWithHeader[H]) UsableSpace() uint16 {
	var h H
	return uint16(bwh.size - utils.GetSize(h))
}

func (bwh *BufferWithHeader[H]) AsSlice() []byte {
	headerPtr := unsafe.Pointer(bwh.header)
	return unsafe.Slice((*byte)(headerPtr), bwh.size)
}
