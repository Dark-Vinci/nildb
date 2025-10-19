package interfaces

type RepPage interface {
	Type() string
	IsOverFlow() bool
	ToBytes() ([]byte, error)
	FromBytes([]byte) error
}

func NewPage(pageSize uint) (RepPage, error) {

}
