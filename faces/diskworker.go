package faces

import "github.com/dark-vinci/nildb/base"

type DiskRequest struct {
	Type       base.DiskOperation
	PageNumber base.PageNumber
	Page       PageHandle
	ResultChan chan DiskResult
}

type DiskResult struct {
	PageNumber base.PageNumber
	Page       PageHandle
	Error      error
}

type DiskWorkerOps interface {
	Write(pageNumber base.PageNumber, page PageHandle) chan DiskResult
	Read(pageNumber base.PageNumber, page PageHandle) chan DiskResult
	Stop()
}
