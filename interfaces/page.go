package interfaces

type PageHandle interface {
	IsOverflow() bool
	IntoBuffer() interface{}
}
