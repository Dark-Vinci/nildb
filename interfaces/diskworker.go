package interfaces

import "github.com/dark-vinci/nildb/base"

type DiskRequest struct {
	Type       base.DiskOperation
	PageNumber base.PageNumber
	Page       RepPage
	ResultChan chan DiskResult
}

type DiskResult struct {
	PageNumber base.PageNumber
	Page       RepPage
	Error      error
}

type DiskWorkerOps interface {
	Write(pageNumber base.PageNumber, page RepPage) chan DiskResult
	Read(pageNumber base.PageNumber, page RepPage) chan DiskResult
	Stop()
}
