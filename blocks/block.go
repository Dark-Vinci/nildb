package blocks

import (
	"fmt"
	"io"

	"github.com/dark-vinci/nildb/interfaces"
)

type Block struct {
	ioOperator interfaces.IOOperator
	blockSize  int
	pageSize   int
}

var _ interfaces.BlockOperations = (*Block)(nil)

func NewBlock(
	ioOperator interfaces.IOOperator,
	blockSize int,
	pageSize int,
) *Block {
	return &Block{
		ioOperator: ioOperator,
		blockSize:  blockSize,
		pageSize:   pageSize,
	}
}

func (b *Block) Write(pageNumber int, buff []byte) error {
	offset := int64(b.pageSize * pageNumber)

	if _, err := b.ioOperator.Seek(offset, io.SeekStart); err != nil {
		fmt.Println("Seek error:", err)
		return err
	}

	if _, err := b.ioOperator.Write(buff); err != nil {
		return err
	}

	return nil
}

func (b *Block) Flush() error {
	return nil
}

func (b *Block) Sync() error {
	if err := b.ioOperator.Sync(); err != nil {
		return err
	}

	return nil
}

func (b *Block) Read(pageNumber int, buff []byte) error {
	var (
		capacity    int
		blockOffset int
		pageOffset  int
	)

	if b.pageSize >= b.blockSize {
		capacity = b.pageSize
		blockOffset = pageNumber * b.pageSize
		pageOffset = 0
	} else {
		offset := (pageNumber * b.pageSize) & ^(b.blockSize - 1)

		capacity = b.blockSize
		blockOffset = offset
		pageOffset = pageNumber*b.pageSize - offset
	}

	_, err := b.ioOperator.Seek(int64(blockOffset), io.SeekStart)
	if err != nil {
		fmt.Println("Error: file cannot be seeked", err)
		return err
	}

	if b.pageSize >= b.blockSize {
		_, err := b.ioOperator.Read(buff)
		if err != nil {
			fmt.Println("Error: file cannot be read", err)
			return err
		}
	}

	block := make([]byte, capacity)
	_, err = b.ioOperator.Read(block)
	if err != nil {
		fmt.Println("Error: file cannot be read", err)
		return err
	}

	copy(buff, block[pageOffset:pageOffset+b.pageSize])

	return nil
}
