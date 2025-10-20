package diskscheduler

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/blocks"
	"github.com/dark-vinci/nildb/constants"
	"github.com/dark-vinci/nildb/interfaces"
)

type DiskWorker struct {
	queue     chan interfaces.DiskRequest
	blockIO   blocks.Block
	lock      sync.RWMutex
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
	pageSize  uint
}

var _ interfaces.DiskWorkerOps = (*DiskWorker)(nil)

func NewDiskWorker(block blocks.Block) *DiskWorker {
	worker := &DiskWorker{
		queue:     make(chan interfaces.DiskRequest, 100),
		blockIO:   block,
		stopChan:  make(chan struct{}),
		lock:      sync.RWMutex{},
		waitGroup: sync.WaitGroup{},
	}

	worker.waitGroup.Add(1)

	go worker.start()

	return worker
}

func (w *DiskWorker) start() {
	defer w.waitGroup.Done()

	var (
		batch  = make([]interfaces.DiskRequest, 0, constants.BatchSize)
		ticker = time.NewTicker(constants.BatchTimeout)
	)

	for {
		select {
		case <-w.stopChan:
			//: Process request before shutting down
			w.processBatch(batch)
			return

		case <-ticker.C:
			if len(batch) > 0 {
				w.processBatch(batch)

				//reset batch
				batch = batch[:0]
			}

		case req := <-w.queue:
			batch = append(batch, req)
			if len(batch) >= constants.BatchSize {
				w.processBatch(batch)

				// reset batch
				batch = batch[:0]

				//reset ticker
				ticker.Reset(constants.BatchTimeout)
			}
		}
	}
}

func (w *DiskWorker) processBatch(batch []interfaces.DiskRequest) {
	// Separate reads and writes
	var (
		reads  []interfaces.DiskRequest
		writes []interfaces.DiskRequest
	)

	for _, req := range batch {
		if req.Type == base.ReadOp {
			reads = append(reads, req)
		} else {
			writes = append(writes, req)
		}
	}

	// Process reads (no reordering needed)
	for _, req := range reads {
		w.processRead(req)
	}

	// Reorder writes by PageNumber
	sort.Slice(writes, func(i, j int) bool {
		return writes[i].PageNumber < writes[j].PageNumber
	})

	for _, req := range writes {
		w.processWrite(req)
	}
}

func (w *DiskWorker) processRead(req interfaces.DiskRequest) {
	data := make([]byte, w.pageSize)

	w.lock.RLock()
	err := w.blockIO.Read(int(req.PageNumber), data)
	w.lock.RUnlock()

	result := interfaces.DiskResult{PageNumber: req.PageNumber, Page: req.Page}

	if err != nil {
		result.Error = fmt.Errorf("page %d not found on disk", req.PageNumber)
	}
	//} else {
	//	if err := req.Page.FromBytes(data); err != nil {
	//		result.Error = fmt.Errorf("failed to deserialize page %d: %v", req.PageNumber, err)
	//	}
	//}

	req.ResultChan <- result
}

func (w *DiskWorker) processWrite(req interfaces.DiskRequest) {
	data, err := req.Page.IntoBuffer()

	if err != nil {
		req.ResultChan <- interfaces.DiskResult{
			PageNumber: req.PageNumber,
			Error:      fmt.Errorf("failed to serialize page %d: %v", req.PageNumber, err),
		}

		return
	}

	pData := data.([]byte)

	w.lock.Lock()
	err = w.blockIO.Write(int(req.PageNumber), pData)
	w.lock.Unlock()

	result := interfaces.DiskResult{PageNumber: req.PageNumber, Page: req.Page}

	if err != nil {
		result.Error = fmt.Errorf("page %d not found on disk", req.PageNumber)
	}

	req.ResultChan <- interfaces.DiskResult{PageNumber: req.PageNumber, Page: req.Page}
}
