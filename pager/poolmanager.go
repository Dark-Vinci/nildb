package pager

import (
	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/interfaces"
	"github.com/dark-vinci/nildb/utils"
	"sync"
)

type Pager struct {
	worker         interfaces.DiskWorkerOps
	cache          interfaces.Cache
	lock           sync.Mutex
	freePages      utils.Uint64Heap
	nextPageNumber base.PageNumber
}
