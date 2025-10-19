package bufferwheader

import (
	"fmt"
	"unsafe"

	"github.com/dark-vinci/nildb/constants"
	"github.com/dark-vinci/nildb/utils"
)

func ForPage[H any](size int) *BufferWithHeader[H] {
	if size < constants.MinPageSize || size > constants.MaxPageSize {
		panic(fmt.Sprintf("INVALID: Page size %v is not between range %v and %v", size, constants.MinPageSize, constants.MaxPageSize))
	}

	return NewBufferWithHeader[H](size)
}

func FromSlice[H any](data []byte) *BufferWithHeader[H] {
	var (
		header     H
		headerSize = utils.GetSize(header)
	)

	if len(data) < headerSize {
		panic(fmt.Sprintf("Attemp to create a BufferWithHeader[%T] from a slice of invalid size of %v", header, len(data)))
	}

	ptr := unsafe.Pointer(&data[0])

	if uintptr(ptr)%constants.CellAlignment != 0 {
		panic(fmt.Sprintf("Attempt to create a BufferWithHeader[%T] from an unaligned pointer %v with %v", header, uintptr(ptr), constants.CellAlignment))
	}

	headerPtr := (*H)(ptr)
	content := data[headerSize:]

	return &BufferWithHeader[H]{
		header:  headerPtr,
		content: content,
		size:    len(data),
	}
}

func Cast[H, T any](b *BufferWithHeader[H]) *BufferWithHeader[T] {
	var t T
	headerSize := utils.GetSize(t)

	if b.size <= headerSize {
		panic(fmt.Sprintf(
			"cannot cast BufferWithHeader[%T] of total size %d to BufferWithHeader[%T] where size of %T is %d",
			*b.header, b.size, t, t, headerSize,
		))
	}

	newHeader := (*T)(unsafe.Pointer(&b.AsSlice()[0]))
	newContent := b.AsSlice()[headerSize:]

	return &BufferWithHeader[T]{
		header:  newHeader,
		content: newContent,
		size:    b.size,
	}
}
