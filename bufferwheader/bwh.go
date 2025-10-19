package bufferwheader

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/dark-vinci/nildb/constants"
)

// BufferWithHeader buffer with header
type BufferWithHeader[H any] struct {
	header  *H
	content []byte
	size    int
}

func getSize[T any](value T) int {
	return int(reflect.TypeOf(value).Size())
}

func Allocate[H any](size int) []byte {
	var (
		header     H
		headerSize = getSize(header)
	)

	if headerSize >= size {
		panic(fmt.Sprintf(
			"An Attempt to allocate BufferWithHeader[%T] of insufficient size: of %T is %d but allocation size is %d",
			header, header, headerSize, size,
		))
	}

	if size%constants.PageAlignment != 0 {
		panic(fmt.Sprintf("Attempt to allocate size: %v that does not match PageAlignment: %v", size, constants.PageAlignment))
	}

	return make([]byte, size)
}

func NewBufferWithHeader[H any](size int) *BufferWithHeader[H] {
	var (
		header      H
		emptyBuffer = Allocate[H](size)
		headerSize  = getSize(header)
		headerPtr   = (*H)(unsafe.Pointer(&emptyBuffer[0]))
	)

	*headerPtr = header

	content := emptyBuffer[headerSize:]

	return &BufferWithHeader[H]{
		header:  headerPtr,
		content: content,
		size:    size,
	}
}
