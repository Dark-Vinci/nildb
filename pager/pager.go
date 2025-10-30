package pager

import (
	"sync"

	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/interfaces"
	"github.com/dark-vinci/nildb/utils"
)

type Pager struct {
	worker         faces.DiskWorkerOps
	cache          faces.Cache
	lock           sync.Mutex
	freePages      utils.Uint64Heap
	nextPageNumber base.PageNumber
}
