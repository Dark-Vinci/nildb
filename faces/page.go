package faces

type PageHandle interface {
	IsOverflow() bool
	IntoBuffer() (interface{}, error)
	FromBuffer([]byte) PageHandle
	Type() string
}
