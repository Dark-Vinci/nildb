package interfaces

type PageHandle interface {
	IsOverflow() bool
	IntoBuffer() (interface{}, error)
	FromBuffer([]byte) PageHandle
}
