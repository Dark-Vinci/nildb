package faces

type BlockOperations interface {
	Sync() error
	Flush() error
	Read(pageNumber int, buff []byte) error
	Write(pageNumber int, buff []byte) error
}
